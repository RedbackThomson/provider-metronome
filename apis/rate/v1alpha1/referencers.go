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

package v1alpha1

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/reference"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	productv1alpha1 "github.com/redbackthomson/provider-metronome/apis/product/v1alpha1"
	ratecardv1alpha1 "github.com/redbackthomson/provider-metronome/apis/ratecard/v1alpha1"
)

// ResolveReferences of this Rate
func (ra *Rate) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, ra)

	var rsp reference.ResolutionResponse
	var err error

	// Resolve spec.forProvider.RateCardID
	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: ra.Spec.ForProvider.RateCardID,
		Reference:    ra.Spec.ForProvider.RateCardRef,
		Selector:     ra.Spec.ForProvider.RateCardSelector,
		To:           reference.To{Managed: &ratecardv1alpha1.RateCard{}, List: &ratecardv1alpha1.RateCardList{}},
		Extract:      RateCardID(),
	})

	if err != nil {
		return errors.Wrap(err, "Spec.ForProvider.RateCardID")
	}

	if rsp.ResolvedValue == "" {
		return errors.New("Spec.ForProvider.RateCardID not yet resolvable")
	}

	ra.Spec.ForProvider.RateCardID = rsp.ResolvedValue
	ra.Spec.ForProvider.RateCardRef = rsp.ResolvedReference

	// Resolve spec.forProvider.ProductID
	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: ra.Spec.ForProvider.ProductID,
		Reference:    ra.Spec.ForProvider.ProductRef,
		Selector:     ra.Spec.ForProvider.ProductSelector,
		To:           reference.To{Managed: &productv1alpha1.Product{}, List: &productv1alpha1.ProductList{}},
		Extract:      ProductID(),
	})

	if err != nil {
		return errors.Wrap(err, "Spec.ForProvider.ProductID")
	}

	if rsp.ResolvedValue == "" {
		return errors.New("Spec.ForProvider.ProductID not yet resolvable")
	}

	ra.Spec.ForProvider.ProductID = rsp.ResolvedValue
	ra.Spec.ForProvider.ProductRef = rsp.ResolvedReference

	return nil
}

// RateCardID extracts info from a kubernetes referenced object
func RateCardID() reference.ExtractValueFn {
	return func(mg resource.Managed) string {
		cr, _ := mg.(*ratecardv1alpha1.RateCard)
		return cr.Status.AtProvider.ID
	}
}

// ProductID extracts info from a kubernetes referenced object
func ProductID() reference.ExtractValueFn {
	return func(mg resource.Managed) string {
		cr, _ := mg.(*productv1alpha1.Product)
		return cr.Status.AtProvider.ID
	}
}
