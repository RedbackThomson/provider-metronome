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
	"encoding/json"
	"fmt"
	"net/http"
)

// EventTypeFilter defines the filter based on event types.
type EventTypeFilter struct {
	InValues    []string `json:"in_values,omitempty"`
	NotInValues []string `json:"not_in_values,omitempty"`
}

// PropertyFilter defines a filter on properties.
type PropertyFilter struct {
	Name        string   `json:"name"`
	Exists      bool     `json:"exists"`
	InValues    []string `json:"in_values,omitempty"`
	NotInValues []string `json:"not_in_values,omitempty"`
}

// CreateBillableMetricRequest represents the request payload for creating a billable metric.
type CreateBillableMetricRequest struct {
	Name            string           `json:"name"`
	EventTypeFilter EventTypeFilter  `json:"event_type_filter"`
	PropertyFilters []PropertyFilter `json:"property_filters"`
	AggregationType string           `json:"aggregation_type"`
	AggregationKey  string           `json:"aggregation_key,omitempty"`
	GroupKeys       [][]string       `json:"group_keys"`
}

// BillableMetric represents the data structure of a billable metric.
type BillableMetric struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	EventTypeFilter EventTypeFilter  `json:"event_type_filter"`
	PropertyFilters []PropertyFilter `json:"property_filters"`
	AggregationType string           `json:"aggregation_type"`
	AggregationKey  string           `json:"aggregation_key,omitempty"`
	GroupKeys       [][]string       `json:"group_keys"`
}

// CreateBillableMetricResponse represents the response for creating a billable metric.
type CreateBillableMetricResponse struct {
	Data BillableMetric `json:"data"`
}

// ListBillableMetricsResponse represents the response for listing billable metrics.
type ListBillableMetricsResponse struct {
	Data     []BillableMetric `json:"data"`
	NextPage *string          `json:"next_page,omitempty"`
}

// UpdateBillableMetricRequest represents the request payload for updating a billable metric.
type UpdateBillableMetricRequest struct {
	Name string `json:"name"`
}

// CreateBillableMetric creates a new billable metric.
func (c *Client) CreateBillableMetric(reqData CreateBillableMetricRequest) (*CreateBillableMetricResponse, error) {
	url := fmt.Sprintf("%s/v1/billable-metrics/create", c.BaseURL)

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := c.newAuthenticatedRequest("POST", url, jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create billable metric: %s", resp.Status)
	}

	var response CreateBillableMetricResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetBillableMetric retrieves a billable metric by ID.
func (c *Client) GetBillableMetric(id string) (*BillableMetric, error) {
	url := fmt.Sprintf("%s/v1/billable-metrics/%s", c.BaseURL, id)

	req, err := c.newAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get billable metric: %s", resp.Status)
	}

	var response struct {
		Data BillableMetric `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}

// ListBillableMetrics retrieves a list of all billable metrics.
func (c *Client) ListBillableMetrics() (*ListBillableMetricsResponse, error) {
	url := fmt.Sprintf("%s/v1/billable-metrics", c.BaseURL)

	req, err := c.newAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list billable metrics: %s", resp.Status)
	}

	var response ListBillableMetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// UpdateBillableMetric updates a billable metric by ID.
func (c *Client) UpdateBillableMetric(id string, reqData UpdateBillableMetricRequest) (*CreateBillableMetricResponse, error) {
	url := fmt.Sprintf("%s/v1/billable-metrics/%s", c.BaseURL, id)

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := c.newAuthenticatedRequest("PUT", url, jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update billable metric: %s", resp.Status)
	}

	var response CreateBillableMetricResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
