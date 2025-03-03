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
	"context"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/client/model"
	"github.com/ONLYOFFICE/onlyoffice-mattermost/server/internal/crypto"
	"github.com/go-resty/resty/v2"
)

var _ OnlyofficeCommandClient = (*onlyofficeCommandClient)(nil)

const (
	OnlyofficeCommandServicePath    string = "/command"
	OnlyofficeCommandServiceVersion string = "version"
)

type OnlyofficeCommandClient interface {
	SendVersion(commandURL string, request model.CommandVersionRequest, timeout time.Duration) (model.CommandVersionResponse, error)
}

type onlyofficeCommandClient struct {
	jwtManager crypto.JwtManager
	client     *resty.Client
}

func NewOnlyofficeCommandClient(jwtManager crypto.JwtManager) OnlyofficeCommandClient {
	return onlyofficeCommandClient{
		jwtManager: jwtManager,
		client:     resty.New(),
	}
}

func (c onlyofficeCommandClient) SendVersion(commandURL string, request model.CommandVersionRequest, timeout time.Duration) (model.CommandVersionResponse, error) {
	var err error

	if len(c.jwtManager.GetKey()) > 0 {
		request.Token, err = c.jwtManager.Sign(request)
	}

	if err != nil {
		return model.CommandVersionResponse{Error: 1, Version: "0.0.0"}, err
	}

	var response model.CommandVersionResponse
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
