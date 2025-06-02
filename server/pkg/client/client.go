/**
 *
 * (c) Copyright Ascensio System SIA 2025
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
	"context"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/pkg/crypto"
	resty "github.com/go-resty/resty/v2"
)

type CommandClient interface {
	SendVersion(
		commandURL string,
		request VersionRequest,
		timeout time.Duration,
	) (VersionResponse, error)
	SendConvert(
		commandURL string,
		request ConvertRequest,
		timeout time.Duration,
	) (ConvertResponse, error)
}

type commandClientImpl struct {
	jwtManager crypto.JwtManager
	client     *resty.Client
}

func New(jwtManager crypto.JwtManager) CommandClient {
	return &commandClientImpl{
		jwtManager: jwtManager,
		client:     resty.New(),
	}
}

func (c *commandClientImpl) SendVersion(
	commandURL string,
	request VersionRequest,
	timeout time.Duration,
) (VersionResponse, error) {
	var response VersionResponse
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err := c.client.R().
		SetBody(request).
		SetResult(&response).
		SetContext(ctx).
		Post(commandURL); err != nil {
		return response, err
	}

	return response, nil
}

func (c *commandClientImpl) SendConvert(
	commandURL string,
	request ConvertRequest,
	timeout time.Duration,
) (ConvertResponse, error) {
	var response ConvertResponse
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err := c.client.R().
		SetBody(request).
		SetResult(&response).
		SetContext(ctx).
		Post(commandURL); err != nil {
		return response, err
	}

	return response, nil
}
