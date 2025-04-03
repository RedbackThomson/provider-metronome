/*
Copyright 2020 The Crossplane Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

type QuantityConversion struct {
	ConversionFactor float64 `json:"conversionFactor"`
	Operation        string  `json:"operation"`
	Name             string  `json:"name,omitempty"`
}

type QuantityRounding struct {
	DecimalPlaces  float64 `json:"decimalPlaces"`
	RoundingMethod string  `json:"roundingMethod"`
}

// ProductParameters represents the request payload for creating a product.
type ProductParameters struct {
	// +optional
	BillableMetricID string `json:"billableMetricId"`
	// +optional
	BillableMetricRef *xpv1.Reference `json:"billableMetricRef,omitempty"`

	// +optional
	BillableMetricSelector *xpv1.Selector `json:"billableMetricSelector,omitempty"`

	Name                 string              `json:"name"`
	Type                 string              `json:"type"`
	CompositeProductIDs  []string            `json:"compositeProductIds,omitempty"`
	CompositeTags        []string            `json:"compositeTags,omitempty"`
	ExcludeFreeUsage     bool                `json:"excludeFreeUsage,omitempty"`
	PresentationGroupKey []string            `json:"presentationGroupKey,omitempty"`
	PricingGroupKey      []string            `json:"pricingGroupKey,omitempty"`
	QuantityConversion   *QuantityConversion `json:"quantityConversion,omitempty"`
	QuantityRounding     *QuantityRounding   `json:"quantityRounding,omitempty"`
	Tags                 []string            `json:"tags,omitempty"`
	StartingAt           string              `json:"startingAt,omitempty"`
}

type ProductDetails struct {
	CreatedAt            string              `json:"createdAt"`
	CreatedBy            string              `json:"createdBy"`
	Name                 string              `json:"name"`
	StartingAt           string              `json:"startingAt,omitempty"`
	CompositeProductIDs  []string            `json:"compositeProductIds,omitempty"`
	CompositeTags        []string            `json:"compositeTags,omitempty"`
	ExcludeFreeUsage     bool                `json:"excludeFreeUsage,omitempty"`
	PresentationGroupKey []string            `json:"presentationGroupKey,omitempty"`
	PricingGroupKey      []string            `json:"pricingGroupKey,omitempty"`
	QuantityConversion   *QuantityConversion `json:"quantityConversion,omitempty"`
	QuantityRounding     *QuantityRounding   `json:"quantityRounding,omitempty"`
	Tags                 []string            `json:"tags,omitempty"`
	BillableMetricID     string              `json:"billableMetricId,omitempty"`
}

// ObservedProduct represents the data structure of a product.
type ObservedProduct struct {
	ID           string            `json:"id"`
	Type         string            `json:"type"`
	Initial      ProductDetails    `json:"initial"`
	Current      ProductDetails    `json:"current"`
	Updates      []ProductDetails  `json:"updates"`
	CustomFields map[string]string `json:"customFields,omitempty"`
	ArchivedAt   string            `json:"archivedAt,omitempty"`
}

// ProductSpec defines the desired state of a Product.
type ProductSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       ProductParameters `json:"forProvider"`
}

// ProductStatus represents the observed state of a Product.
type ProductStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ObservedProduct `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// Product represents a Metronome Billable Metric resource
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,metronome}
type Product struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProductSpec   `json:"spec"`
	Status ProductStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProductList contains a list of Product
type ProductList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Product `json:"items"`
}
