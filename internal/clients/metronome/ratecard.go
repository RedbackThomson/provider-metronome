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

	"github.com/redbackthomson/provider-metronome/apis/ratecard/v1alpha1"
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

type UpdateRateCardResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

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
		return nil, ErrInvalidName
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
		return nil, ErrInvalidName
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
		return nil, errors.New("failed to update rate card: " + resp.Status)
	}

	var response UpdateRateCardResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// RateCardConverter helps to convert Metronome client types to api types
// of this provider and vise-versa From & To shall both be defined for each type
// conversion, to prevent divergence from Metronome client Types
// goverter:converter
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
// goverter:extend ExtV1JSONToRuntimeRawExtension
// goverter:enum:unknown @ignore
// goverter:struct:comment // +k8s:deepcopy-gen=false
// goverter:output:file ./zz_generated.ratecard.conversion.go
// +k8s:deepcopy-gen=false
type RateCardConverter interface {
	FromRateCardSpec(in *v1alpha1.RateCardParameters) *CreateRateCardRequest
	ToRateCardSpec(in *CreateRateCardRequest) *v1alpha1.RateCardParameters

	FromRateCard(in *RateCard) *v1alpha1.ObservedRateCard
	ToRateCard(in *v1alpha1.ObservedRateCard) *RateCard

	// goverter:ignoreMissing
	FromRateCardToParameters(in *RateCard) *v1alpha1.RateCardParameters
}
