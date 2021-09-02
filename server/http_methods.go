/**
 *
 * (c) Copyright Ascensio System SIA 2021
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

package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type HTTPClient struct {
	client http.Client
}

type Header struct {
	Key   string
	Value string
}

//TODO: Rebuild this function
func (httpClient HTTPClient) PostRequest(url string, requestBody interface{}, headers []Header,
	responseBody interface {
		Succeeded()
		Failed()
	}) {
	body := &requestBody
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	req, err := http.NewRequest("POST", url, buf)

	if err != nil {
		responseBody.Failed()
		return
	}

	if len(headers) > 0 {
		for _, header := range headers {
			req.Header.Add(header.Key, header.Value)
		}
	}

	res, err := httpClient.client.Do(req)

	if err != nil {
		responseBody.Failed()
		return
	}

	if res.StatusCode < 300 {
		responseBody.Succeeded()
	} else {
		responseBody.Failed()
	}

	defer httpClient.client.CloseIdleConnections()

	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(responseBody)
}

func (httpClient HTTPClient) GetRequest(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Response error: ", err.Error())
		return nil, err
	}
	defer httpClient.client.CloseIdleConnections()
	return resp, nil
}
