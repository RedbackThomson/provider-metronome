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
)

type GetRatesRequest struct {
	RateCardID string         `json:"rate_card_id"`
	At         string         `json:"at"`
	Selectors  []RateSelector `json:"selectors"`
}

type RateSelector struct {
	ProductID                 string            `json:"product_id"`
	PartialPricingGroupValues map[string]string `json:"partial_pricing_group_values"`
}

type GetRatesResponse struct {
	// Response structure to be defined later
}

type AddRatesRequest struct {
	RateCardID string `json:"rate_card_id"`
	Rates      []Rate `json:"rates"`
}

type Rate struct {
	ProductID          string            `json:"product_id"`
	StartingAt         string            `json:"starting_at"`
	Entitled           bool              `json:"entitled"`
	RateType           string            `json:"rate_type"`
	Price              float64           `json:"price"`
	PricingGroupValues map[string]string `json:"pricing_group_values"`
}

type AddRatesResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (c *Client) GetRates(reqData GetRatesRequest) (*GetRatesResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/rate-cards/getRates", c.BaseURL)
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := c.newAuthenticatedRequest("POST", url, jsonData)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get rates: " + resp.Status)
	}

	var response GetRatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) AddRates(reqData AddRatesRequest) (*AddRatesResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/rate-cards/addRates", c.BaseURL)
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := c.newAuthenticatedRequest("POST", url, jsonData)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to add rates: " + resp.Status)
	}

	var response AddRatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
