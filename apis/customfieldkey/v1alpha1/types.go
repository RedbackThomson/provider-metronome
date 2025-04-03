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

// CustomFieldKeyParameters represents the request payload for creating a custom field key.
type CustomFieldKeyParameters struct {
	EnforceUniqueness bool   `json:"enforceUniqueness"`
	Entity            string `json:"entity"`
	Key               string `json:"key"`
}

// ObservedCustomFieldKey represents the data structure of a custom field key.
type ObservedCustomFieldKey struct {
	EnforceUniqueness bool   `json:"enforceUniqueness"`
	Entity            string `json:"entity"`
	Key               string `json:"key"`
}

// CustomFieldKeySpec defines the desired state of a CustomFieldKey.
type CustomFieldKeySpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       CustomFieldKeyParameters `json:"forProvider"`
}

// CustomFieldKeyStatus represents the observed state of a CustomFieldKey.
type CustomFieldKeyStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ObservedCustomFieldKey `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// CustomFieldKey represents a Metronome Custom field key resource
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,metronome}
type CustomFieldKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomFieldKeySpec   `json:"spec"`
	Status CustomFieldKeyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CustomFieldKeyList contains a list of CustomFieldKey
type CustomFieldKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomFieldKey `json:"items"`
}
