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

	billablemetricv1alpha1 "github.com/redbackthomson/provider-metronome/apis/billablemetric/v1alpha1"
)

// ResolveReferences of this Product
func (pr *Product) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, pr)

	var rsp reference.ResolutionResponse
	var err error

	// Resolve spec.forProvider.BillableMetricID
	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: pr.Spec.ForProvider.BillableMetricID,
		Reference:    pr.Spec.ForProvider.BillableMetricRef,
		Selector:     pr.Spec.ForProvider.BillableMetricSelector,
		To:           reference.To{Managed: &billablemetricv1alpha1.BillableMetric{}, List: &billablemetricv1alpha1.BillableMetricList{}},
		Extract:      BillableMetricID(),
	})

	if err != nil {
		return errors.Wrap(err, "Spec.ForProvider.BillableMetricID")
	}

	if rsp.ResolvedValue == "" {
		return errors.New("Spec.ForProvider.BillableMetricID not yet resolvable")
	}

	pr.Spec.ForProvider.BillableMetricID = rsp.ResolvedValue
	pr.Spec.ForProvider.BillableMetricRef = rsp.ResolvedReference

	return nil
}

// BillableMetricID extracts info from a kubernetes referenced object
func BillableMetricID() reference.ExtractValueFn {
	return func(mg resource.Managed) string {
		cr, _ := mg.(*billablemetricv1alpha1.BillableMetric)
		return cr.Status.AtProvider.ID
	}
}
