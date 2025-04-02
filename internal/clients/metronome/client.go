/*
Copyright 2025 RedbackThomson.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metronome

import (
	"bytes"
	"net/http"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
)

type Client struct {
	logger logging.Logger

	baseURL    string
	authToken  string
	httpClient *http.Client
}

func (c *Client) newAuthenticatedRequest(method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func New(log logging.Logger, baseURL, authToken string) *Client {
	return &Client{
		logger:     log,
		baseURL:    baseURL,
		authToken:  authToken,
		httpClient: &http.Client{},
	}
}
