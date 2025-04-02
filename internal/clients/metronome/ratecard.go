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
	ErrRateCardInvalidName = errors.New("invalid rate card name")
)

type GetRateCardRequest struct {
	ID string `json:"id"`
}

type GetRateCardResponse struct {
	Data RateCard `json:"data"`
}

type RateCard struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	FiatCreditType FiatCreditType    `json:"fiat_credit_type,omitempty"`
	CreatedAt      string            `json:"created_at"`
	CreatedBy      string            `json:"created_by"`
	Aliases        []RateCardAlias   `json:"aliases,omitempty"`
	CustomFields   map[string]string `json:"custom_fields,omitempty"`
}

type CreateRateCardRequest struct {
	Name                  string                 `json:"name"`
	Description           string                 `json:"description,omitempty"`
	FiatCreditTypeID      string                 `json:"fiat_credit_type_id,omitempty"`
	CreditTypeConversions []CreditTypeConversion `json:"credit_type_conversions,omitempty"`
	Aliases               []RateCardAlias        `json:"aliases,omitempty"`
	CustomFields          map[string]string      `json:"custom_fields,omitempty"`
}

type CreditTypeConversion struct {
	CustomCreditTypeID  string `json:"custom_credit_type_id"`
	FiatPerCustomCredit string `json:"fiat_per_custom_credit"`
}

type CreateRateCardResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

type UpdateRateCardRequest struct {
	RateCardID  string `json:"rate_card_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateRateCardResponse DataID

type FiatCreditType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RateCardAlias struct {
	Name string `json:"name"`
}

func (c *Client) GetRateCard(reqData GetRateCardRequest) (*GetRateCardResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/rate-cards/get", c.baseURL)

	if !IsUUID(reqData.ID) {
		return nil, ErrRateCardInvalidName
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
			return nil, errors.Wrap(c, "failed to get rate card")
		}
		return nil, errors.New("failed to get rate card: " + resp.Status)
	}

	var response GetRateCardResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) CreateRateCard(reqData CreateRateCardRequest) (*CreateRateCardResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/rate-cards/create", c.baseURL)
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
			return nil, errors.Wrap(c, "failed to create rate card")
		}
		return nil, errors.New("failed to create rate card: " + resp.Status)
	}

	var response CreateRateCardResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) UpdateRateCard(reqData UpdateRateCardRequest) (*UpdateRateCardResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/rate-cards/update", c.baseURL)

	if !IsUUID(reqData.RateCardID) {
		return nil, ErrRateCardInvalidName
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
			return nil, errors.Wrap(c, "failed to update rate card")
		}
		return nil, errors.New("failed to update rate card: " + resp.Status)
	}

	var response UpdateRateCardResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
