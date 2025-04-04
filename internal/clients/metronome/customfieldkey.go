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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrCustomFieldKeyInvalidName = errors.New("invalid custom field key name")
)

type CustomFieldKey struct {
	EnforceUniqueness bool   `json:"enforce_uniqueness"`
	Entity            string `json:"entity"`
	Key               string `json:"key"`
}

// CreateCustomFieldKeyRequest represents the request payload for creating a custom field key.
type CreateCustomFieldKeyRequest CustomFieldKey

// ListCustomFieldKeysRequest represents the request payload for listing custom field keys.
type ListCustomFieldKeysRequest struct {
	Entities []string `json:"entities,omitempty"`
}

// ListCustomFieldKeysResponse represents the response for listing custom field keys.
type ListCustomFieldKeysResponse struct {
	Data     []CustomFieldKey `json:"data"`
	NextPage string           `json:"next_page,omitempty"`
}

// DeleteCustomFieldKeyRequest represents the request payload to delete a custom field key.
type DeleteCustomFieldKeyRequest struct {
	Entity string `json:"entity"`
	Key    string `json:"key"`
}

// CreateCustomFieldKey creates a new custom field key.
func (c *Client) CreateCustomFieldKey(ctx context.Context, reqData CreateCustomFieldKeyRequest) error {
	url := fmt.Sprintf("%s/v1/customFields/addKey", c.baseURL)

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := c.newAuthenticatedRequest(ctx, "POST", url, jsonData)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() // nolint:errcheck // Read-only stream

	if resp.StatusCode != http.StatusOK {
		if c := ParseClientError(resp.Body); c != nil {
			return fmt.Errorf("failed to create custom field key: %s", c.Message)
		}
		return fmt.Errorf("failed to create custom field key: %s", resp.Status)
	}

	return nil
}

// ListCustomFieldKeys retrieves a list of all custom field keys.
func (c *Client) ListCustomFieldKeys(ctx context.Context, reqData ListCustomFieldKeysRequest, nextPage string) (*ListCustomFieldKeysResponse, error) {
	url := fmt.Sprintf("%s/v1/customFields/listKeys", c.baseURL)

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := c.newAuthenticatedRequest(ctx, "POST", url, jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	if nextPage != "" {
		q.Add("next_page", nextPage)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() // nolint:errcheck // Read-only stream

	if resp.StatusCode != http.StatusOK {
		if c := ParseClientError(resp.Body); c != nil {
			return nil, fmt.Errorf("failed to list custom field keys: %s", c.Message)
		}
		return nil, fmt.Errorf("failed to list custom field keys: %s", resp.Status)
	}

	var response ListCustomFieldKeysResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// DeleteCustomFieldKey deletes a custom field key by ID.
func (c *Client) DeleteCustomFieldKey(ctx context.Context, reqData DeleteCustomFieldKeyRequest) error {
	url := fmt.Sprintf("%s/v1/customFields/removeKey", c.baseURL)

	// Prepare the request payload
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Create a new authenticated request
	req, err := c.newAuthenticatedRequest(ctx, "POST", url, jsonData)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() // nolint:errcheck // Read-only stream

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		if c := ParseClientError(resp.Body); c != nil {
			return fmt.Errorf("failed to delete custom field key: %s", c.Message)
		}
		return fmt.Errorf("failed to delete custom field key: %s", resp.Status)
	}

	return nil
}
