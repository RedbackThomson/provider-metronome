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

	"github.com/pkg/errors"
)

var (
	ErrContractInvalidName = errors.New("invalid contract name")
	ErrCustomerInvalidName = errors.New("invalid contract name")
)

type GetContractRequest struct {
	CustomerID string `json:"customer_id"`
	ContractID string `json:"contract_id"`
}

type ProductIdentifier struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Commit struct {
	ID                   string            `json:"id"`
	Type                 string            `json:"type"`
	Name                 string            `json:"name"`
	Description          string            `json:"description"`
	Product              ProductIdentifier `json:"product"`
	RolloverFraction     float64           `json:"rollover_fraction"`
	ApplicableProductIDs []string          `json:"applicable_product_ids"`
}

type Override struct {
	ID         string            `json:"id"`
	Product    ProductIdentifier `json:"product"`
	StartingAt string            `json:"starting_at"`
	Type       string            `json:"type"`
	Multiplier float64           `json:"multiplier"`
}

type Contract struct {
	ID                  string            `json:"id"`
	CustomerID          string            `json:"customer_id"`
	RateCardID          string            `json:"rate_card_id"`
	StartingAt          string            `json:"starting_at"`
	NetPaymentTermsDays int               `json:"net_payment_terms_days"`
	EndingBefore        string            `json:"ending_before"`
	CreatedAt           string            `json:"created_at"`
	CreatedBy           string            `json:"created_by"`
	CustomFields        map[string]string `json:"custom_fields"`
	Commits             []Commit          `json:"commits"`
	Overrides           []Override        `json:"overrides"`
}

type GetContractResponse struct {
	Data Contract `json:"data"`
}

type ListContractsRequest struct {
	CustomerID string `json:"customer_id"`
}

type ListContractsResponse struct {
	Data []Contract `json:"data"`
}

func (c *Client) GetContract(reqData GetContractRequest) (*GetContractResponse, error) {
	url := fmt.Sprintf("%s/v2/contracts/get", c.baseURL)

	if !IsUUID(reqData.ContractID) {
		return nil, ErrContractInvalidName
	}
	if !IsUUID(reqData.CustomerID) {
		return nil, ErrCustomerInvalidName
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := c.newAuthenticatedRequest("POST", url, jsonData)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if c := ParseClientError(resp.Body); c != nil {
			return nil, errors.Wrap(c, "failed to get contract")
		}
		return nil, errors.New("failed to get contract: " + resp.Status)
	}

	var response GetContractResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) ListContracts(reqData ListContractsRequest) (*ListContractsResponse, error) {
	url := fmt.Sprintf("%s/v2/contracts/list", c.baseURL)
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := c.newAuthenticatedRequest("POST", url, jsonData)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if c := ParseClientError(resp.Body); c != nil {
			return nil, errors.Wrap(c, "failed to list contracts")
		}
		return nil, errors.New("failed to list contracts: " + resp.Status)
	}

	var response ListContractsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
