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

package billablemetric

import (
	"context"
	"slices"

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

	"github.com/redbackthomson/provider-metronome/apis/billablemetric/v1alpha1"
	metronomev1alpha1 "github.com/redbackthomson/provider-metronome/apis/v1alpha1"
	metronomeClient "github.com/redbackthomson/provider-metronome/internal/clients/metronome"
	"github.com/redbackthomson/provider-metronome/internal/connector"
	"github.com/redbackthomson/provider-metronome/internal/converters"
)

const (
	errNotBillableMetric     = "managed resource is not a BillableMetric custom resource"
	errGetBillableMetric     = "failed to get billable metric"
	errCreateBillableMetric  = "failed to create billable metric"
	errArchiveBillableMetric = "failed to archive billable metric"
)

// Setup adds a controller that reconciles BillableMetric managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, baseUrl string) error {
	name := managed.ControllerName(v1alpha1.BillableMetricGroupKind)

	reconcilerOptions := []managed.ReconcilerOption{
		managed.WithExternalConnecter(
			&connector.Connector[*v1alpha1.BillableMetric, *metronomeExternal]{
				Logger:               o.Logger,
				Client:               mgr.GetClient(),
				Usage:                resource.NewProviderConfigUsageTracker(mgr.GetClient(), &metronomev1alpha1.ProviderConfigUsage{}),
				BaseURL:              baseUrl,
				NewMetronomeClientFn: metronomeClient.New,
				NewExternalClientFn: func(log logging.Logger, client *metronomeClient.Client) *metronomeExternal {
					return &metronomeExternal{
						logger:    o.Logger,
						metronome: client.BillableMetric(),
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
		mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, &v1alpha1.BillableMetricList{}, o.MetricOptions.PollStateMetricInterval)); err != nil {
		return err
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.BillableMetricGroupVersionKind),
		reconcilerOptions...,
	)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.BillableMetric{}).
		WithOptions(o.ForControllerRuntime()).
		Complete(r)
}

type metronomeExternal struct {
	logger    logging.Logger
	metronome metronomeClient.BillableMetricClient
}

func (e *metronomeExternal) Disconnect(ctx context.Context) error {
	return nil
}

func (e *metronomeExternal) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.BillableMetric)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotBillableMetric)
	}

	e.logger.Debug("Observing")

	id := meta.GetExternalName(cr)
	if id == "" {
		return managed.ExternalObservation{}, nil
	}

	res, err := e.metronome.GetBillableMetric(ctx, id)
	if err != nil {
		// the external name isn't valid
		if errors.Is(err, metronomeClient.ErrBillableMetricInvalidName) {
			return managed.ExternalObservation{ResourceExists: false}, nil
		}
		return managed.ExternalObservation{}, errors.Wrap(err, errGetBillableMetric)
	}

	if res == nil {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	metric := &res.Data

	if metric.ArchivedAt != "" {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	converter := &converters.BillableMetricConverterImpl{}
	cr.Status.AtProvider = *converter.FromBillableMetric(metric)
	cr.SetConditions(xpv1.Available())

	upToDate, diff := isUpToDate(cr, metric)

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: upToDate,
		Diff:             diff,
	}, nil
}

func (e *metronomeExternal) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.BillableMetric)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotBillableMetric)
	}

	e.logger.Debug("Creating")

	converter := &converters.BillableMetricConverterImpl{}
	req := converter.FromBillableMetricSpec(&cr.Spec.ForProvider)

	res, err := e.metronome.CreateBillableMetric(ctx, *req)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateBillableMetric)
	}
	if res.Data.ID == "" {
		return managed.ExternalCreation{}, errors.New("billable metric ID is missing")
	}

	meta.SetExternalName(cr, res.Data.ID)

	return managed.ExternalCreation{}, nil
}

func (e *metronomeExternal) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	_, ok := mg.(*v1alpha1.BillableMetric)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotBillableMetric)
	}

	e.logger.Debug("Updating")

	return managed.ExternalUpdate{}, errors.New("updating a billable metric is not supported")
}

func (e *metronomeExternal) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.BillableMetric)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotBillableMetric)
	}

	e.logger.Debug("Deleting")

	id := meta.GetExternalName(cr)
	if id == "" {
		return managed.ExternalDelete{}, nil
	}

	_, err := e.metronome.ArchiveBillableMetric(ctx, id)
	if err != nil {
		if errors.Is(err, metronomeClient.ErrBillableMetricAlreadyArchived) {
			return managed.ExternalDelete{}, nil
		}
		return managed.ExternalDelete{}, errors.Wrap(err, errArchiveBillableMetric)
	}

	return managed.ExternalDelete{}, nil
}

func isUpToDate(cr *v1alpha1.BillableMetric, metric *metronomeClient.BillableMetric) (bool, string) {
	spec := cr.Spec.ForProvider.DeepCopy()

	converter := &converters.BillableMetricConverterImpl{}
	params := converter.FromBillableMetricToParameters(metric)

	sortPropertyFilter := func(a, b v1alpha1.PropertyFilter) int {
		if a.Name < b.Name {
			return 1
		}
		if a.Name == b.Name {
			return 0
		}
		return -1
	}

	slices.SortFunc(spec.PropertyFilters, sortPropertyFilter)
	slices.SortFunc(params.PropertyFilters, sortPropertyFilter)

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
	}

	return cmp.Equal(spec, params, opts...), cmp.Diff(spec, params, opts...)
}
