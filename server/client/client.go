/**
 *
 * (c) Copyright Ascensio System SIA 2023
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
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/client/model"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/crypto"
)

var _ OnlyofficeCommandClient = (*onlyofficeCommandClient)(nil)

const (
	OnlyofficeCommandServicePath    string = "/coauthoring/CommandService.ashx"
	OnlyofficeCommandServiceVersion string = "version"
)

var failedVersionResponse = model.CommandVersionResponse{Error: 1, Version: "0.0.0"}

type OnlyofficeCommandClient interface {
	SendVersion(commandURL string, request model.CommandVersionRequest, timeout time.Duration) (model.CommandVersionResponse, error)
}

type onlyofficeCommandClient struct {
	jwtManager crypto.JwtManager
	client     http.Client
}

func NewOnlyofficeCommandClient(insecure bool, jwtManager crypto.JwtManager) OnlyofficeCommandClient {
	return onlyofficeCommandClient{
		jwtManager: jwtManager,
		client: http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: insecure},
			},
		},
	}
}

func (c onlyofficeCommandClient) SendVersion(commandURL string, request model.CommandVersionRequest, timeout time.Duration) (model.CommandVersionResponse, error) {
	var err error

	if len(c.jwtManager.GetKey()) > 0 {
		request.Token, err = c.jwtManager.Sign(request)
	}

	if err != nil {
		return failedVersionResponse, err
	}

	buf, err := json.Marshal(request)

	if err != nil {
		return failedVersionResponse, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, commandURL, bytes.NewBuffer(buf))

	if err != nil {
		return failedVersionResponse, err
	}

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return failedVersionResponse, err
	}

	response := model.CommandVersionResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return failedVersionResponse, err
	}

	return response, nil
}
