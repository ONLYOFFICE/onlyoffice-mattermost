/**
 *
 * (c) Copyright Ascensio System SIA 2026
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package health

import (
	"sync"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/bot"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/client"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/configuration"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	cron "github.com/robfig/cron/v3"
)

const (
	CronSchedule           = "*/30 * * * *"
	onlyofficeLoggerPrefix = "[ONLYOFFICE] "
)

type HealthChecker interface {
	Start()
	Stop()
	IsHealthy() bool
}

type healthChecker struct {
	configuration *configuration.Configuration

	api           plugin.API
	bot           bot.Bot
	commandClient client.CommandClient

	jwtManager crypto.JwtManager

	cron        *cron.Cron
	isHealthy   bool
	healthMutex sync.RWMutex
	isRunning   bool
	runMutex    sync.Mutex
}

func NewHealthChecker(
	configuration *configuration.Configuration,
	api plugin.API,
	bot bot.Bot,
	commandClient client.CommandClient,
	jwtManager crypto.JwtManager,
) HealthChecker {
	return &healthChecker{
		configuration: configuration,
		api:           api,
		bot:           bot,
		commandClient: commandClient,
		jwtManager:    jwtManager,
		isHealthy:     true,
		cron:          cron.New(),
	}
}

func (h *healthChecker) setHealthy(healthy bool) {
	h.healthMutex.Lock()
	defer h.healthMutex.Unlock()
	h.isHealthy = healthy
}

func (h *healthChecker) notifyAdmins(message string) {
	page := 0
	perPage := 100

	for {
		users, err := h.api.GetUsers(&model.UserGetOptions{
			Page:    page,
			PerPage: perPage,
		})

		if err != nil {
			h.api.LogError(onlyofficeLoggerPrefix + "Could not get users for admin notification: " + err.Message)
			return
		}

		if len(users) == 0 {
			break
		}

		for _, user := range users {
			if user.IsSystemAdmin() {
				h.bot.BotCreateDM(message, user.Id)
			}
		}

		if len(users) < perPage {
			break
		}

		page++
	}
}

func (h *healthChecker) publishHealthStatus(healthy bool) {
	event := map[string]any{
		"healthy": healthy,
		"time":    time.Now().Unix(),
	}

	h.api.PublishWebSocketEvent(
		"health_status",
		event,
		&model.WebsocketBroadcast{},
	)
}

func (h *healthChecker) handleUnhealthy() {
	h.setHealthy(false)

	message := `⚠️ ONLYOFFICE Plugin Configuration Issue

The ONLYOFFICE plugin configuration is invalid. Editor buttons have been disabled for all users.

Please check your plugin settings:

• Document Server Address
• JWT Secret
• JWT Header
• JWT Prefix

Or enable the demo mode if you want to test the plugin.`

	h.notifyAdmins(message)
	h.publishHealthStatus(false)
}

func (h *healthChecker) createVersionToken() (string, error) {
	vreq := client.VersionRequest{Command: "version"}
	vreq.IssuedAt = jwt.NewNumericDate(time.Now())
	vreq.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Minute))
	return h.jwtManager.Sign([]byte(h.configuration.DESJwt), vreq)
}

func (h *healthChecker) sendVersionRequest(token string) (client.VersionResponse, error) {
	resp, err := h.commandClient.SendVersion(
		h.configuration.DESAddress+client.OnlyofficeCommandServicePath+"?shardkey="+uuid.New().String(),
		client.VersionRequest{
			Command: "version",
			Token:   token,
		},
		4*time.Second,
	)

	return resp, err
}

func (h *healthChecker) checkHealth() {
	if err := h.configuration.IsValid(); err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "Configuration validation failed: " + err.Error())
		h.handleUnhealthy()
		return
	}

	token, err := h.createVersionToken()
	if err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "Could not create version token: " + err.Error())
		h.handleUnhealthy()
		return
	}

	resp, err := h.sendVersionRequest(token)
	if err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "Could not reach document server: " + err.Error())
		h.handleUnhealthy()
		return
	}

	if resp.Error != 0 {
		h.api.LogError(onlyofficeLoggerPrefix + "Document server returned error code: " + string(rune(resp.Error)))
		h.handleUnhealthy()
		return
	}

	wasUnhealthy := !h.IsHealthy()
	h.setHealthy(true)

	if wasUnhealthy {
		h.api.LogInfo(onlyofficeLoggerPrefix + "Document server health check passed - service recovered")
		h.publishHealthStatus(true)
	} else {
		h.api.LogDebug(onlyofficeLoggerPrefix + "Document server health check passed")
	}
}

func (h *healthChecker) Start() {
	h.runMutex.Lock()
	defer h.runMutex.Unlock()

	if h.isRunning {
		return
	}

	h.api.LogInfo(onlyofficeLoggerPrefix + "Starting periodic health checks")
	_, err := h.cron.AddFunc(CronSchedule, func() {
		h.checkHealth()
	})

	if err != nil {
		h.api.LogError(onlyofficeLoggerPrefix + "Failed to add cron job: " + err.Error())
		return
	}

	h.cron.Start()
	h.isRunning = true

	go h.checkHealth()
}

func (h *healthChecker) Stop() {
	h.runMutex.Lock()
	defer h.runMutex.Unlock()

	if !h.isRunning {
		return
	}

	h.api.LogInfo(onlyofficeLoggerPrefix + "Stopping periodic health checks")

	if h.cron != nil {
		ctx := h.cron.Stop()
		<-ctx.Done()
	}

	h.isRunning = false
}

func (h *healthChecker) IsHealthy() bool {
	h.healthMutex.RLock()
	defer h.healthMutex.RUnlock()
	return h.isHealthy
}
