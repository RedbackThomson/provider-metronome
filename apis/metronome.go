/*
Copyright 2020 RedbackThomson.

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

// Package apis contains Kubernetes API for the Metronome provider.
package apis

import (
	"k8s.io/apimachinery/pkg/runtime"

	billablemetricv1alpha1 "github.com/redbackthomson/provider-metronome/apis/billablemetric/v1alpha1"
	customfieldkeyv1alpha1 "github.com/redbackthomson/provider-metronome/apis/customfieldkey/v1alpha1"
	productv1alpha1 "github.com/redbackthomson/provider-metronome/apis/product/v1alpha1"
	ratev1alpha1 "github.com/redbackthomson/provider-metronome/apis/rate/v1alpha1"
	ratecardv1alpha1 "github.com/redbackthomson/provider-metronome/apis/ratecard/v1alpha1"
	metronomev1alpha1 "github.com/redbackthomson/provider-metronome/apis/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes,
		metronomev1alpha1.SchemeBuilder.AddToScheme,
		billablemetricv1alpha1.SchemeBuilder.AddToScheme,
		customfieldkeyv1alpha1.SchemeBuilder.AddToScheme,
		productv1alpha1.SchemeBuilder.AddToScheme,
		ratecardv1alpha1.SchemeBuilder.AddToScheme,
		ratev1alpha1.SchemeBuilder.AddToScheme,
	)
}

// AddToSchemes may be used to add all resources defined in the project to a Scheme
var AddToSchemes runtime.SchemeBuilder

// AddToScheme adds all Resources to the Scheme
func AddToScheme(s *runtime.Scheme) error {
	return AddToSchemes.AddToScheme(s)
}
