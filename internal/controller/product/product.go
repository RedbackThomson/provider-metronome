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

package product

import (
	"context"
	"sort"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/feature"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/statemetrics"

	"github.com/redbackthomson/provider-metronome/apis/product/v1alpha1"
	metronomev1alpha1 "github.com/redbackthomson/provider-metronome/apis/v1alpha1"
	metronomeClient "github.com/redbackthomson/provider-metronome/internal/clients/metronome"
	"github.com/redbackthomson/provider-metronome/internal/connector"
	"github.com/redbackthomson/provider-metronome/internal/converters"
)

const (
	errNotProduct     = "managed resource is not a Product custom resource"
	errGetProduct     = "failed to get product"
	errCreateProduct  = "failed to create product"
	errUpdateProduct  = "failed to update product"
	errArchiveProduct = "failed to archive product"
	errNoID           = "product does not have ID"
	errNoStartingAt   = "forProvider.startingAt is required for updates"
)

// Setup adds a controller that reconciles Product managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, baseUrl string) error {
	name := managed.ControllerName(v1alpha1.ProductGroupKind)

	reconcilerOptions := []managed.ReconcilerOption{
		managed.WithExternalConnecter(
			&connector.Connector[*v1alpha1.Product, *metronomeExternal]{
				Logger:               o.Logger,
				Client:               mgr.GetClient(),
				Usage:                resource.NewProviderConfigUsageTracker(mgr.GetClient(), &metronomev1alpha1.ProviderConfigUsage{}),
				BaseURL:              baseUrl,
				NewMetronomeClientFn: metronomeClient.New,
				NewExternalClientFn: func(log logging.Logger, client *metronomeClient.Client) *metronomeExternal {
					return &metronomeExternal{
						logger:    o.Logger,
						metronome: client.Product(),
					}
				},
			}),
		managed.WithPollInterval(o.PollInterval),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithMetricRecorder(o.MetricOptions.MRMetrics),
	}

	if o.Features.Enabled(feature.EnableBetaManagementPolicies) {
		reconcilerOptions = append(reconcilerOptions, managed.WithManagementPolicies())
	}

	if err := mgr.Add(statemetrics.NewMRStateRecorder(
		mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, &v1alpha1.ProductList{}, o.MetricOptions.PollStateMetricInterval)); err != nil {
		return err
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.ProductGroupVersionKind),
		reconcilerOptions...,
	)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.Product{}).
		WithOptions(o.ForControllerRuntime()).
		Complete(r)
}

type metronomeExternal struct {
	logger    logging.Logger
	metronome metronomeClient.ProductClient
}

func (e *metronomeExternal) Disconnect(ctx context.Context) error {
	return nil
}

func (e *metronomeExternal) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Product)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotProduct)
	}

	e.logger.Debug("Observing")

	id := meta.GetExternalName(cr)
	if id == "" {
		return managed.ExternalObservation{}, nil
	}

	res, err := e.metronome.GetProduct(ctx, metronomeClient.GetProductRequest{
		ID: id,
	})
	if err != nil {
		// the external name isn't valid
		if errors.Is(err, metronomeClient.ErrProductInvalidName) {
			return managed.ExternalObservation{ResourceExists: false}, nil
		}
		return managed.ExternalObservation{}, errors.Wrap(err, errGetProduct)
	}

	if res == nil {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	pr := res.Data

	if pr.ArchivedAt != "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	card := &res.Data

	converter := &converters.ProductConverterImpl{}
	cr.Status.AtProvider = *converter.FromProduct(card)
	cr.SetConditions(xpv1.Available())

	upToDate, diff := isUpToDate(cr, card)

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: upToDate,
		Diff:             diff,
	}, nil
}

func (e *metronomeExternal) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Product)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotProduct)
	}

	e.logger.Debug("Creating")

	converter := &converters.ProductConverterImpl{}
	req := converter.FromProductSpec(&cr.Spec.ForProvider)

	res, err := e.metronome.CreateProduct(ctx, *req)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateProduct)
	}
	if res.Data.ID == "" {
		return managed.ExternalCreation{}, errors.New("product ID is missing")
	}

	meta.SetExternalName(cr, res.Data.ID)

	return managed.ExternalCreation{}, nil
}

func (e *metronomeExternal) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Product)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotProduct)
	}

	e.logger.Debug("Updating")

	id := meta.GetExternalName(cr)

	if id == "" {
		return managed.ExternalUpdate{}, errors.New(errNoID)
	}

	if cr.Spec.ForProvider.StartingAt == "" {
		return managed.ExternalUpdate{}, errors.New(errNoStartingAt)
	}

	converter := &converters.ProductConverterImpl{}
	req := converter.ToProductUpdate(&cr.Spec.ForProvider)
	req.ProductID = id

	res, err := e.metronome.UpdateProduct(ctx, *req)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateProduct)
	}
	if res.Data.ID == "" {
		return managed.ExternalUpdate{}, errors.New("product ID is missing")
	}

	return managed.ExternalUpdate{}, nil
}

func (e *metronomeExternal) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.Product)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotProduct)
	}

	e.logger.Debug("Deleting")

	id := meta.GetExternalName(cr)
	if id == "" {
		return managed.ExternalDelete{}, nil
	}

	_, err := e.metronome.ArchiveProduct(ctx, metronomeClient.ArchiveProductRequest{
		ProductID: id,
	})
	if err != nil {
		if errors.Is(err, metronomeClient.ErrProductAlreadyArchived) {
			return managed.ExternalDelete{}, nil
		}
		return managed.ExternalDelete{}, errors.Wrap(err, errArchiveProduct)
	}

	return managed.ExternalDelete{}, nil
}

func isUpToDate(cr *v1alpha1.Product, metric *metronomeClient.Product) (bool, string) {
	spec := cr.Spec.ForProvider.DeepCopy()

	converter := &converters.ProductConverterImpl{}
	params := converter.FromProductToParameters(metric)

	caseInsensitiveComparer := cmp.Comparer(strings.EqualFold)

	spec.BillableMetricRef = nil
	spec.BillableMetricSelector = nil

	sort.Strings(spec.CompositeProductIDs)
	sort.Strings(spec.CompositeTags)
	sort.Strings(spec.PresentationGroupKey)
	sort.Strings(spec.PricingGroupKey)
	sort.Strings(spec.Tags)

	sort.Strings(params.CompositeProductIDs)
	sort.Strings(params.CompositeTags)
	sort.Strings(params.PresentationGroupKey)
	sort.Strings(params.PricingGroupKey)
	sort.Strings(params.Tags)

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmp.FilterPath(func(p cmp.Path) bool {
			return p.String() == "Type"
		}, caseInsensitiveComparer),
		cmpopts.IgnoreFields(v1alpha1.ProductParameters{},
			"StartingAt",
		),
	}

	return cmp.Equal(spec, params, opts...), cmp.Diff(spec, params, opts...)
}
