package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type HTTPClient struct {
	client http.Client
}

//TODO: Rebuild this function
func (httpClient HTTPClient) PostRequest(url string, requestBody interface{}, responseBody interface {
	Connected()
	CheckResponse() error
}) {
	body := &requestBody
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	req, _ := http.NewRequest("POST", url, buf)

	res, _ := httpClient.client.Do(req)

	if res.StatusCode < 300 {
		responseBody.Connected()
	}

	defer httpClient.client.CloseIdleConnections()

	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(responseBody)
}

func (httpClient HTTPClient) GetRequest(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Response error: ", err.Error())
	}
	return resp.Body
}
