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

// EventTypeFilter defines the filter based on event types.
type EventTypeFilter struct {
	InValues    []string `json:"in_values,omitempty"`
	NotInValues []string `json:"not_in_values,omitempty"`
}

// PropertyFilter defines a filter on properties.
type PropertyFilter struct {
	Name        string   `json:"name"`
	Exists      *bool    `json:"exists,omitempty"`
	InValues    []string `json:"in_values,omitempty"`
	NotInValues []string `json:"not_in_values,omitempty"`
}

type AggregationType string

const (
	AggregationCount  = "count"
	AggregationLatest = "latest"
	AggregationMax    = "max"
	AggregationSum    = "sum"
	AggregationUnique = "unique"
)

// BillableMetricParameters represents the request payload for creating a billable metric.
type BillableMetricParameters struct {
	Name            string            `json:"name"`
	AggregationType AggregationType   `json:"aggregation_type"`
	AggregationKey  string            `json:"aggregation_key"`
	EventTypeFilter EventTypeFilter   `json:"event_type_filter"`
	PropertyFilters []PropertyFilter  `json:"property_filters"`
	GroupKeys       [][]string        `json:"group_keys"`
	CustomFields    map[string]string `json:"custom_fields,omitempty"`
	SQL             string            `json:"sql,omitempty"`
}

// ObservedBillableMetric represents the data structure of a billable metric.
type ObservedBillableMetric struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// +kubebuilder:validation:Enum=count;latest;max;sum;unique
	AggregationType AggregationType   `json:"aggregation_type"`
	AggregationKey  string            `json:"aggregation_key,omitempty"`
	EventTypeFilter EventTypeFilter   `json:"event_type_filter"`
	PropertyFilters []PropertyFilter  `json:"property_filters"`
	GroupKeys       [][]string        `json:"group_keys"`
	CustomFields    map[string]string `json:"custom_fields,omitempty"`
	SQL             string            `json:"sql,omitempty"`
	ArchivedAt      string            `json:"archived_at,omitempty"`
}

// BillableMetricSpec defines the desired state of a Release.
type BillableMetricSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       BillableMetricParameters `json:"forProvider"`
}

// BillableMetricStatus represents the observed state of a Release.
type BillableMetricStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ObservedBillableMetric `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// BillableMetric represents a Metronome Billable Metric resource
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,metronome}
type BillableMetric struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BillableMetricSpec   `json:"spec"`
	Status BillableMetricStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BillableMetricList contains a list of Release
type BillableMetricList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BillableMetric `json:"items"`
}
