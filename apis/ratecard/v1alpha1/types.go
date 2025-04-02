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

type RateCardAlias struct {
	Name string `json:"name"`
}

type CreditTypeConversion struct {
	CustomCreditTypeID  string `json:"customCreditTypeId"`
	FiatPerCustomCredit string `json:"fiatPerCustomCredit"`
}

type FiatCreditType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RateCardParameters represents the request payload for creating a rate card.
type RateCardParameters struct {
	Name                  string                 `json:"name"`
	Description           string                 `json:"description,omitempty"`
	FiatCreditTypeID      string                 `json:"fiatCreditTypeId,omitempty"`
	CreditTypeConversions []CreditTypeConversion `json:"creditTypeConversions,omitempty"`
	Aliases               []RateCardAlias        `json:"aliases,omitempty"`
	CustomFields          map[string]string      `json:"customFields,omitempty"`
}

// ObservedRateCard represents the data structure of a rate card.
type ObservedRateCard struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	FiatCreditType FiatCreditType    `json:"fiatCreditType,omitempty"`
	CreatedAt      string            `json:"createdAt"`
	CreatedBy      string            `json:"createdBy"`
	Aliases        []RateCardAlias   `json:"aliases,omitempty"`
	CustomFields   map[string]string `json:"customFields,omitempty"`
}

// RateCardSpec defines the desired state of a Release.
type RateCardSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       RateCardParameters `json:"forProvider"`
}

// RateCardStatus represents the observed state of a Release.
type RateCardStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ObservedRateCard `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// RateCard represents a Metronome Rate Card resource
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,metronome}
type RateCard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RateCardSpec   `json:"spec"`
	Status RateCardStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RateCardList contains a list of Release
type RateCardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RateCard `json:"items"`
}
