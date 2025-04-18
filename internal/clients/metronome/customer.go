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
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type CustomerClient interface {
	CreateCustomer(ctx context.Context, reqData CreateCustomerRequest) (*CreateCustomerResponse, error)
	GetCustomer(ctx context.Context, customerID string) (*GetCustomerResponse, error)
	UpdateCustomerAliases(ctx context.Context, customerID string, reqData UpdateAliasesRequest) error
	ListCustomers(ctx context.Context) (*ListCustomersResponse, error)
}

type CustomerClientImpl struct {
	Client *Client
}

var _ (CustomerClient) = (*CustomerClientImpl)(nil)

type CreateCustomerRequest struct {
	IngestAliases                 []string                       `json:"ingest_aliases"`
	Name                          string                         `json:"name"`
	BillingProviderConfigurations []BillingProviderConfiguration `json:"customer_billing_provider_configurations"`
}

type CreateCustomerResponse struct {
	Data CreateCustomerData `json:"data"`
}

type CreateCustomerData struct {
	ID            string   `json:"id"`
	ExternalID    string   `json:"external_id"`
	IngestAliases []string `json:"ingest_aliases"`
	Name          string   `json:"name"`
}

type GetCustomerResponse struct {
	Data GetCustomerData `json:"data"`
}

type GetCustomerData struct {
	ID             string            `json:"id"`
	ExternalID     string            `json:"external_id"`
	IngestAliases  []string          `json:"ingest_aliases"`
	Name           string            `json:"name"`
	CustomerConfig CustomerConfig    `json:"customer_config"`
	CustomFields   map[string]string `json:"custom_fields"`
}

type ListCustomersResponse struct {
	Data     []GetCustomerData `json:"data"`
	NextPage *string           `json:"next_page"`
}

type BillingProviderConfiguration struct {
	BillingProvider string               `json:"billing_provider"`
	DeliveryMethod  string               `json:"delivery_method"`
	Configuration   BillingConfiguration `json:"configuration"`
}

type BillingConfiguration struct {
	StripeCustomerID       string `json:"stripe_customer_id"`
	StripeCollectionMethod string `json:"stripe_collection_method"`
}

type CustomerConfig struct {
	SalesforceAccountID string `json:"salesforce_account_id"`
}

type UpdateAliasesRequest struct {
	IngestAliases []string `json:"ingest_aliases"`
}

func (c *CustomerClientImpl) CreateCustomer(ctx context.Context, reqData CreateCustomerRequest) (*CreateCustomerResponse, error) {
	url := fmt.Sprintf("%s/v1/customers", c.Client.baseURL)
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := c.Client.newAuthenticatedRequest(ctx, "POST", url, jsonData)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint:errcheck // Read-only stream

	if resp.StatusCode != http.StatusOK {
		if c := ParseClientError(resp.Body); c != nil {
			return nil, errors.Wrap(c, "failed to create customer")
		}
		return nil, errors.New("failed to create customer: " + resp.Status)
	}

	var response CreateCustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *CustomerClientImpl) GetCustomer(ctx context.Context, customerID string) (*GetCustomerResponse, error) {
	url := fmt.Sprintf("%s/v1/customers/%s", c.Client.baseURL, customerID)

	req, err := c.Client.newAuthenticatedRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint:errcheck // Read-only stream

	if resp.StatusCode != http.StatusOK {
		if c := ParseClientError(resp.Body); c != nil {
			return nil, errors.Wrap(c, "failed to get customer")
		}
		return nil, errors.New("failed to get customer: " + resp.Status)
	}

	var response GetCustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *CustomerClientImpl) UpdateCustomerAliases(ctx context.Context, customerID string, reqData UpdateAliasesRequest) error {
	url := fmt.Sprintf("%s/v1/customers/%s/setIngestAliases", c.Client.baseURL, customerID)
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	req, err := c.Client.newAuthenticatedRequest(ctx, "POST", url, jsonData)
	if err != nil {
		return err
	}

	resp, err := c.Client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint:errcheck // Read-only stream

	if resp.StatusCode != http.StatusOK {
		if c := ParseClientError(resp.Body); c != nil {
			return errors.Wrap(c, "failed to update customer aliases")
		}
		return errors.New("failed to update customer aliases: " + resp.Status)
	}

	return nil
}

func (c *CustomerClientImpl) ListCustomers(ctx context.Context) (*ListCustomersResponse, error) {
	url := fmt.Sprintf("%s/v1/customers", c.Client.baseURL)
	req, err := c.Client.newAuthenticatedRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint:errcheck // Read-only stream

	if resp.StatusCode != http.StatusOK {
		if c := ParseClientError(resp.Body); c != nil {
			return nil, errors.Wrap(c, "failed to list customers")
		}
		return nil, errors.New("failed to list customers: " + resp.Status)
	}

	var response ListCustomersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
