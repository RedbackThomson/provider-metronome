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

package converters

import (
	"github.com/redbackthomson/provider-metronome/apis/rate/v1alpha1"
	"github.com/redbackthomson/provider-metronome/internal/clients/metronome"
)

// RateConverter helps to convert Metronome client types to api types
// of this provider and vise-versa From & To shall both be defined for each type
// conversion, to prevent divergence from Metronome client Types
// goverter:converter
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
// goverter:enum:unknown @ignore
// goverter:struct:comment // +k8s:deepcopy-gen=false
// goverter:output:file ./zz_generated.rate.conversion.go
// +k8s:deepcopy-gen=false
type RateConverter interface {
	FromRateSpec(in *v1alpha1.RateParameters) *metronome.AddRateRequest

	// goverter:ignore RateCardRef RateCardSelector ProductRef ProductSelector
	ToRateSpec(in *metronome.AddRateRequest) *v1alpha1.RateParameters

	FromRate(in *metronome.Rate) *v1alpha1.ObservedRate
	ToRate(in *v1alpha1.ObservedRate) *metronome.Rate

	// goverter:ignoreMissing
	// goverter:ignore RateCardID
	// goverter:map Details.RateType RateType
	// goverter:map Details.IsProrated IsProrated
	// goverter:map Details.Price Price
	// goverter:map Details.PricingGroupValues PricingGroupValues
	// goverter:map Details.Quantity Quantity
	// goverter:map Details.Tiers Tiers
	// goverter:map Details.UseListPrices UseListPrices
	// goverter:map Details.CreditType.ID CreditTypeID
	FromRateToParameters(in *metronome.Rate) *v1alpha1.RateParameters
}
