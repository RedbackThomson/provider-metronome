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

var (
	ErrRateInvalidName = errors.New("invalid rate name")
)

type GetRatesRequest struct {
	RateCardID string         `json:"rate_card_id"`
	At         string         `json:"at"`
	Selectors  []RateSelector `json:"selectors"`
}

type RateSelector struct {
	ProductID                 string            `json:"product_id,omitempty"`
	PartialPricingGroupValues map[string]string `json:"partial_pricing_group_values,omitempty"`
	PricingGroupValues        map[string]string `json:"pricing_group_values,omitempty"`
	ProductTags               []string          `json:"product_tags,omitempty"`
}

type GetRatesResponse struct {
	Data     []Rate `json:"data"`
	NextPage string `json:"next_page"`
}

type CreditType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AddRatesRequest struct {
	RateCardID string        `json:"rate_card_id"`
	Rates      []RateDetails `json:"rates"`
}

type AddRateRequest struct {
	Entitled           bool              `json:"entitled"`
	ProductID          string            `json:"product_id"`
	RateCardID         string            `json:"rate_card_id"`
	RateType           string            `json:"rate_type"`
	StartingAt         string            `json:"starting_at"`
	CommitRate         *CommitRate       `json:"commit_rate,omitempty"`
	CreditTypeID       string            `json:"credit_type_id,omitempty"`
	EndingBefore       string            `json:"ending_before,omitempty"`
	IsProrated         bool              `json:"is_prorated,omitempty"`
	Price              float64           `json:"price,omitempty"`
	PricingGroupValues map[string]string `json:"pricing_group_values,omitempty"`
	Quantity           float64           `json:"quantity,omitempty"`
	Tiers              []Tier            `json:"tiers,omitempty"`
	UseListPrices      bool              `json:"use_list_prices,omitempty"`
}

type AddRateResponse struct {
	Data struct {
		RateType string  `json:"rate_type"`
		Price    float64 `json:"price"`
	} `json:"data"`
}

type Tier struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type CommitRate struct {
	RateType string  `json:"rate_type"`
	Price    float64 `json:"price"`
	Tiers    []Tier  `json:"tiers"`
}

type Rate struct {
	Entitled           bool              `json:"entitled"`
	ProductCustomField map[string]string `json:"product_custom_fields"`
	ProductID          string            `json:"product_id"`
	ProductName        string            `json:"product_name"`
	ProductTags        []string          `json:"product_tags"`
	Details            RateDetails       `json:"rate"`
	StartingAt         string            `json:"starting_at"`
	CommitRate         *CommitRate       `json:"commit_rate,omitempty"`
	EndingBefore       string            `json:"ending_before,omitempty"`
	PricingGroupValues map[string]string `json:"pricing_group_values,omitempty"`
}

type RateDetails struct {
	RateType           string            `json:"rate_type"`
	CreditType         CreditType        `json:"credit_type,omitempty"`
	IsProrated         bool              `json:"is_prorated,omitempty"`
	Price              float64           `json:"price,omitempty"`
	PricingGroupValues map[string]string `json:"pricing_group_values,omitempty"`
	Quantity           float64           `json:"quantity,omitempty"`
	Tiers              []Tier            `json:"tiers,omitempty"`
	UseListPrices      bool              `json:"use_list_prices,omitempty"`
}

type AddRatesResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (c *Client) GetRates(ctx context.Context, reqData GetRatesRequest, nextPage string) (*GetRatesResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/rate-cards/getRates", c.baseURL)

	if !IsUUID(reqData.RateCardID) {
		return nil, ErrRateCardInvalidName
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := c.newAuthenticatedRequest(ctx, "POST", url, jsonData)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if nextPage != "" {
		q.Add("next_page", nextPage)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint:errcheck // Read-only stream

	if resp.StatusCode != http.StatusOK {
		if c := ParseClientError(resp.Body); c != nil {
			return nil, errors.Wrap(c, "failed to get rates")
		}
		return nil, errors.New("failed to get rates: " + resp.Status)
	}

	var response GetRatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) AddRate(ctx context.Context, reqData AddRateRequest) (*AddRateResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/rate-cards/addRate", c.baseURL)

	if !IsUUID(reqData.RateCardID) {
		return nil, ErrRateInvalidName
	}
	if !IsUUID(reqData.ProductID) {
		return nil, ErrProductInvalidName
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := c.newAuthenticatedRequest(ctx, "POST", url, jsonData)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint:errcheck // Read-only stream

	if resp.StatusCode != http.StatusOK {
		if c := ParseClientError(resp.Body); c != nil {
			return nil, errors.Wrap(c, "failed to add rate")
		}
		return nil, errors.New("failed to add rate: " + resp.Status)
	}

	var response AddRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
