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

type Tier struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type CommitRate struct {
	RateType string  `json:"rateType"`
	Price    float64 `json:"price"`
	Tiers    []Tier  `json:"tiers"`
}

// RateParameters represents the request payload for creating a rate card.
type RateParameters struct {
	RateCardID string `json:"rateCardId,omitempty"`

	// +optional
	RateCardRef *xpv1.Reference `json:"rateCardRef,omitempty"`

	// +optional
	RateCardSelector *xpv1.Selector `json:"rateCardSelector,omitempty"`

	ProductID  string `json:"productId"`
	StartingAt string `json:"startingAt"`
	Entitled   bool   `json:"entitled"`
	RateType   string `json:"rateType"`

	// Price is the default price. For FLAT and SUBSCRIPTION rateType, this
	// must be >=0 and the unit is **CENTS**. For PERCENTAGE rateType, this is
	// a decimal fraction, e.g. use 0.1 for 10%; this must be >=0 and <=1.
	Price              float64           `json:"price,omitempty"`
	PricingGroupValues map[string]string `json:"pricingGroupValues,omitempty"`
	CommitRate         *CommitRate       `json:"commitRate,omitempty"`
	CreditTypeID       string            `json:"creditTypeId,omitempty"`
	EndingBefore       string            `json:"endingBefore,omitempty"`
	IsProrated         bool              `json:"isProrated,omitempty"`
	Quantity           float64           `json:"quantity,omitempty"`
	Tiers              []Tier            `json:"tiers,omitempty"`
	UseListPrices      bool              `json:"useListPrices,omitempty"`
}

type CreditType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RateDetails struct {
	RateType           string            `json:"rateType"`
	CreditType         CreditType        `json:"creditType,omitempty"`
	IsProrated         bool              `json:"isProrated,omitempty"`
	Price              float64           `json:"price,omitempty"`
	PricingGroupValues map[string]string `json:"pricingGroupValues,omitempty"`
	Quantity           float64           `json:"quantity,omitempty"`
	Tiers              []Tier            `json:"tiers,omitempty"`
	UseListPrices      bool              `json:"useListPrices,omitempty"`
}

// ObservedRate represents the data structure of a rate card.
type ObservedRate struct {
	Entitled           bool              `json:"entitled"`
	ProductCustomField map[string]string `json:"productCustomFields"`
	ProductID          string            `json:"productId"`
	ProductName        string            `json:"productName"`
	ProductTags        []string          `json:"productTags"`
	Details            RateDetails       `json:"rate"`
	StartingAt         string            `json:"startingAt"`
	CommitRate         CommitRate        `json:"commitRate,omitempty"`
	EndingBefore       string            `json:"endingBefore,omitempty"`
	PricingGroupValues map[string]string `json:"pricingGroupValues,omitempty"`
}

// RateSpec defines the desired state of a Release.
type RateSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       RateParameters `json:"forProvider"`
}

// RateStatus represents the observed state of a Release.
type RateStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ObservedRate `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// Rate represents a Metronome Rate resource
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,metronome}
type Rate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RateSpec   `json:"spec"`
	Status RateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RateList contains a list of Release
type RateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rate `json:"items"`
}
