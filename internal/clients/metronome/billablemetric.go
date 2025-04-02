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
	"errors"
	"fmt"
	"net/http"

	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/google/uuid"
	"github.com/redbackthomson/provider-metronome/apis/billablemetric/v1alpha1"
)

var (
	ErrInvalidName     = errors.New("invalid billable metric name")
	ErrAlreadyArchived = errors.New("billable metric already archived")
)

const (
	errBillableMetricAlreadyArchived = "Billable metric already archived"
)

// EventTypeFilter defines the filter based on event types.
type EventTypeFilter struct {
	InValues    []string `json:"in_values,omitempty"`
	NotInValues []string `json:"not_in_values,omitempty"`
}

// PropertyFilter defines a filter on properties.
type PropertyFilter struct {
	Name        string   `json:"name"`
	Exists      *bool    `json:"exists,omitempty"`
	InValues    []string `json:"in_values,omitempty"`
	NotInValues []string `json:"not_in_values,omitempty"`
}

type AggregationType string

const (
	AggregationCount  = "count"
	AggregationLatest = "latest"
	AggregationMax    = "max"
	AggregationSum    = "sum"
	AggregationUnique = "unique"
)

// CreateBillableMetricRequest represents the request payload for creating a billable metric.
type CreateBillableMetricRequest struct {
	Name            string            `json:"name"`
	AggregationType AggregationType   `json:"aggregation_type"`
	AggregationKey  string            `json:"aggregation_key"`
	EventTypeFilter EventTypeFilter   `json:"event_type_filter"`
	PropertyFilters []PropertyFilter  `json:"property_filters"`
	GroupKeys       [][]string        `json:"group_keys"`
	CustomFields    map[string]string `json:"custom_fields,omitempty"`
	SQL             string            `json:"sql,omitempty"`
}

// BillableMetric represents the data structure of a billable metric.
type BillableMetric struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// +kubebuilder:validation:Enum=count;latest;max;sum;unique
	AggregationType AggregationType   `json:"aggregation_type"`
	AggregationKey  string            `json:"aggregation_key,omitempty"`
	EventTypeFilter EventTypeFilter   `json:"event_type_filter"`
	PropertyFilters []PropertyFilter  `json:"property_filters"`
	GroupKeys       [][]string        `json:"group_keys"`
	CustomFields    map[string]string `json:"custom_fields,omitempty"`
	SQL             string            `json:"sql,omitempty"`
	ArchivedAt      string            `json:"archived_at,omitempty"`
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

// ArchiveBillableMetricRequest represents the request payload to archive a billable metric.
type ArchiveBillableMetricRequest struct {
	ID string `json:"id"` // ID of the billable metric to archive
}

// ArchiveBillableMetricResponse represents the response for archiving a billable metric.
type ArchiveBillableMetricResponse struct {
	Data struct {
		ID string `json:"id"` // ID of the archived billable metric
	} `json:"data"`
}

// CreateBillableMetric creates a new billable metric.
func (c *Client) CreateBillableMetric(reqData CreateBillableMetricRequest) (*CreateBillableMetricResponse, error) {
	url := fmt.Sprintf("%s/v1/billable-metrics/create", c.baseURL)

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := c.newAuthenticatedRequest("POST", url, jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
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
	url := fmt.Sprintf("%s/v1/billable-metrics/%s", c.baseURL, id)

	if !IsValidBillableMetricName(id) {
		return nil, ErrInvalidName
	}

	req, err := c.newAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
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
	url := fmt.Sprintf("%s/v1/billable-metrics", c.baseURL)

	req, err := c.newAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
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
	url := fmt.Sprintf("%s/v1/billable-metrics/%s", c.baseURL, id)

	if !IsValidBillableMetricName(id) {
		return nil, ErrInvalidName
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := c.newAuthenticatedRequest("PUT", url, jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
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

// ArchiveBillableMetric archives a billable metric by ID.
func (c *Client) ArchiveBillableMetric(id string) (*ArchiveBillableMetricResponse, error) {
	url := fmt.Sprintf("%s/v1/billable-metrics/archive", c.baseURL)

	if !IsValidBillableMetricName(id) {
		return nil, ErrInvalidName
	}

	// Prepare the request payload
	reqData := ArchiveBillableMetricRequest{ID: id}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Create a new authenticated request
	req, err := c.newAuthenticatedRequest("POST", url, jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		if IsClientError(resp.Body, errBillableMetricAlreadyArchived) {
			return nil, ErrAlreadyArchived
		}
		return nil, fmt.Errorf("failed to archive billable metric: %s", resp.Status)
	}

	// Decode the response data
	var response ArchiveBillableMetricResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func IsValidBillableMetricName(s string) bool {
	return uuid.Validate(s) == nil
}

// BillableMetricConverter helps to convert Metronome client types to api types
// of this provider and vise-versa From & To shall both be defined for each type
// conversion, to prevent
// divergence from Metronome client Types
// goverter:converter
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
// goverter:extend ExtV1JSONToRuntimeRawExtension
// goverter:enum:unknown @ignore
// goverter:struct:comment // +k8s:deepcopy-gen=false
// goverter:output:file ./zz_generated.billablemetric.conversion.go
// +k8s:deepcopy-gen=false
type BillableMetricConverter interface {
	FromBillableMetricSpec(in *v1alpha1.BillableMetricParameters) *CreateBillableMetricRequest
	ToBillableMetricSpec(in *CreateBillableMetricRequest) *v1alpha1.BillableMetricParameters

	FromBillableMetric(in *BillableMetric) *v1alpha1.ObservedBillableMetric
	ToBillableMetric(in *v1alpha1.ObservedBillableMetric) *BillableMetric

	// goverter:ignoreMissing
	FromBillableMetricToParameters(in *BillableMetric) *v1alpha1.BillableMetricParameters
}

// ExtV1JSONToRuntimeRawExtension converts an extv1.JSON into a
// *runtime.RawExtension.
func ExtV1JSONToRuntimeRawExtension(in extv1.JSON) *runtime.RawExtension {
	return &runtime.RawExtension{
		Raw: in.Raw,
	}
}
