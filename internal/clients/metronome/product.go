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
	ErrProductInvalidName     = errors.New("invalid product name")
	ErrProductAlreadyArchived = errors.New("product already archived")
)

const (
	errProductAlreadyArchived = "Product already archived"
)

type ArchiveProductRequest struct {
	ProductID string `json:"product_id"`
}

type ArchiveProductResponse DataID

type CreateProductRequest struct {
	Name                 string              `json:"name"`
	Type                 string              `json:"type"`
	BillableMetricID     string              `json:"billable_metric_id,omitempty"`
	CompositeProductIDs  []string            `json:"composite_product_ids,omitempty"`
	CompositeTags        []string            `json:"composite_tags,omitempty"`
	ExcludeFreeUsage     bool                `json:"exclude_free_usage,omitempty"`
	PresentationGroupKey []string            `json:"presentation_group_key,omitempty"`
	PricingGroupKey      []string            `json:"pricing_group_key,omitempty"`
	QuantityConversion   *QuantityConversion `json:"quantity_conversion,omitempty"`
	QuantityRounding     *QuantityRounding   `json:"quantity_rounding,omitempty"`
	Tags                 []string            `json:"tags,omitempty"`
}

type CreateProductResponse DataID

type GetProductRequest struct {
	ID string `json:"id"`
}

type ListProductsRequest struct {
	ArchiveFilter string `json:"archive_filter"`
}

type UpdateProductRequest struct {
	ProductID            string              `json:"product_id"`
	StartingAt           string              `json:"starting_at"`
	BillableMetricID     string              `json:"billable_metric_id,omitempty"`
	CompositeProductIDs  []string            `json:"composite_product_ids,omitempty"`
	CompositeTags        []string            `json:"composite_tags,omitempty"`
	ExcludeFreeUsage     bool                `json:"exclude_free_usage,omitempty"`
	Name                 string              `json:"name,omitempty"`
	PresentationGroupKey []string            `json:"presentation_group_key,omitempty"`
	PricingGroupKey      []string            `json:"pricing_group_key,omitempty"`
	QuantityConversion   *QuantityConversion `json:"quantity_conversion,omitempty"`
	QuantityRounding     *QuantityRounding   `json:"quantity_rounding,omitempty"`
	Tags                 []string            `json:"tags,omitempty"`
}

type UpdateProductResponse DataID

type GetProductResponse struct {
	Data Product `json:"data"`
}

type Product struct {
	ID           string            `json:"id"`
	Type         string            `json:"type"`
	Initial      ProductDetails    `json:"initial"`
	Current      ProductDetails    `json:"current"`
	Updates      []ProductDetails  `json:"updates"`
	CustomFields map[string]string `json:"custom_fields"`
	ArchivedAt   string            `json:"archived_at"`
}

type QuantityConversion struct {
	ConversionFactor float64 `json:"conversion_factor"`
	Operation        string  `json:"operation"`
	Name             string  `json:"name"`
}

type QuantityRounding struct {
	DecimalPlaces  float64 `json:"decimal_places"`
	RoundingMethod string  `json:"rounding_method"`
}

type ProductDetails struct {
	Name                 string              `json:"name"`
	BillableMetricID     string              `json:"billable_metric_id,omitempty"`
	CompositeProductIDs  []string            `json:"composite_product_ids"`
	CompositeTags        []string            `json:"composite_tags"`
	ExcludeFreeUsage     bool                `json:"exclude_free_usage,omitempty"`
	PresentationGroupKey []string            `json:"presentation_group_key"`
	PricingGroupKey      []string            `json:"pricing_group_key"`
	QuantityConversion   *QuantityConversion `json:"quantity_conversion,omitempty"`
	QuantityRounding     *QuantityRounding   `json:"quantity_rounding"`
	Tags                 []string            `json:"tags"`
	StartingAt           string              `json:"starting_at"`
	CreatedAt            string              `json:"created_at"`
	CreatedBy            string              `json:"created_by"`
}

type ListProductsResponse struct {
	Data     []Product `json:"data"`
	NextPage *string   `json:"next_page"`
}

func (c *Client) ArchiveProduct(reqData ArchiveProductRequest) (*ArchiveProductResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/products/archive", c.baseURL)

	if !IsUUID(reqData.ProductID) {
		return nil, ErrProductInvalidName
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
			if c.Message == errProductAlreadyArchived {
				return nil, ErrProductAlreadyArchived
			}
			return nil, errors.Wrap(c, "failed to archive product")
		}
		return nil, errors.New("failed to archive product: " + resp.Status)
	}

	var response ArchiveProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) CreateProduct(reqData CreateProductRequest) (*CreateProductResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/products/create", c.baseURL)
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
			return nil, errors.Wrap(c, "failed to create product")
		}
		return nil, errors.New("failed to create product: " + resp.Status)
	}

	var response CreateProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetProduct(reqData GetProductRequest) (*GetProductResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/products/get", c.baseURL)

	if !IsUUID(reqData.ID) {
		return nil, ErrProductInvalidName
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
			return nil, errors.Wrap(c, "failed to get product")
		}
		return nil, errors.New("failed to get product: " + resp.Status)
	}

	var response GetProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) ListProduct(reqData ListProductsRequest) (*ListProductsResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/products/list", c.baseURL)
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
			return nil, errors.Wrap(c, "failed to list product")
		}
		return nil, errors.New("failed to list product: " + resp.Status)
	}

	var response ListProductsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) UpdateProduct(reqData UpdateProductRequest) (*UpdateProductResponse, error) {
	url := fmt.Sprintf("%s/v1/contract-pricing/products/update", c.baseURL)

	if !IsUUID(reqData.ProductID) {
		return nil, ErrProductInvalidName
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
			return nil, errors.Wrap(c, "failed to update product")
		}
		return nil, errors.New("failed to update product: " + resp.Status)
	}

	var response UpdateProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
