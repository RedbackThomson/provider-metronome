/*
Copyright 2021 The Crossplane Authors.

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

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/pkg/controller"

	"github.com/redbackthomson/provider-metronome/internal/controller/billablemetric"
	"github.com/redbackthomson/provider-metronome/internal/controller/config"
	"github.com/redbackthomson/provider-metronome/internal/controller/customfieldkey"
	"github.com/redbackthomson/provider-metronome/internal/controller/product"
	"github.com/redbackthomson/provider-metronome/internal/controller/rate"
	"github.com/redbackthomson/provider-metronome/internal/controller/ratecard"
)

// Setup creates all Template controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options, baseUrl string) error {
	if err := config.Setup(mgr, o); err != nil {
		return err
	}
	if err := billablemetric.Setup(mgr, o, baseUrl); err != nil {
		return err
	}
	if err := customfieldkey.Setup(mgr, o, baseUrl); err != nil {
		return err
	}
	if err := product.Setup(mgr, o, baseUrl); err != nil {
		return err
	}
	if err := rate.Setup(mgr, o, baseUrl); err != nil {
		return err
	}
	if err := ratecard.Setup(mgr, o, baseUrl); err != nil {
		return err
	}
	return nil
}
