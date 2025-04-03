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
	"github.com/redbackthomson/provider-metronome/apis/customfieldkey/v1alpha1"
	"github.com/redbackthomson/provider-metronome/internal/clients/metronome"
)

// CustomFieldKeyConverter helps to convert Metronome client types to api types
// of this provider and vise-versa From & To shall both be defined for each type
// conversion, to prevent divergence from Metronome client Types
// goverter:converter
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
// goverter:enum:unknown @ignore
// goverter:struct:comment // +k8s:deepcopy-gen=false
// goverter:output:file ./zz_generated.customfieldkey.conversion.go
// +k8s:deepcopy-gen=false
type CustomFieldKeyConverter interface {
	FromCustomFieldKeySpec(in *v1alpha1.CustomFieldKeyParameters) *metronome.CreateCustomFieldKeyRequest
	ToCustomFieldKeySpec(in *metronome.CreateCustomFieldKeyRequest) *v1alpha1.CustomFieldKeyParameters

	FromCustomFieldKey(in *metronome.CustomFieldKey) *v1alpha1.ObservedCustomFieldKey
	ToCustomFieldKey(in *v1alpha1.ObservedCustomFieldKey) *metronome.CustomFieldKey

	// goverter:ignoreMissing
	FromCustomFieldKeyToParameters(in *metronome.CustomFieldKey) *v1alpha1.CustomFieldKeyParameters
}
