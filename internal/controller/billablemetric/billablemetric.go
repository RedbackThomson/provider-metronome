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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

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
	"github.com/redbackthomson/provider-metronome/internal/converters"
)

const (
	errNotBillableMetric         = "managed resource is not a BillableMetric custom resource"
	errProviderConfigNotSet      = "provider config is not set"
	errGetProviderConfig         = "cannot get provider config"
	errGetCreds                  = "failed to create credentials from provider config"
	errFailedToTrackUsage        = "cannot track provider config usage"
	errFailedToGetBillableMetric = "failed to get billable metric"
	errCreateBillableMetric      = "failed to create billable metric"
	errArchiveBillableMetric     = "failed to archive billable metric"
)

// Setup adds a controller that reconciles BillableMetric managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, baseUrl string) error {
	name := managed.ControllerName(v1alpha1.BillableMetricGroupKind)

	reconcilerOptions := []managed.ReconcilerOption{
		managed.WithExternalConnecter(&connector{
			client:               mgr.GetClient(),
			logger:               o.Logger,
			usage:                resource.NewProviderConfigUsageTracker(mgr.GetClient(), &metronomev1alpha1.ProviderConfigUsage{}),
			newMetronomeClientFn: metronomeClient.New,
			baseURL:              baseUrl,
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

type connector struct {
	baseURL string
	logger  logging.Logger
	client  client.Client
	usage   resource.Tracker

	newMetronomeClientFn func(log logging.Logger, baseURL, authToken string) *metronomeClient.Client
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) { //nolint:gocyclo
	cr, ok := mg.(*v1alpha1.BillableMetric)
	if !ok {
		return nil, errors.New(errNotBillableMetric)
	}
	l := c.logger.WithValues("request", cr.Name)

	l.Debug("Connecting")

	pc := &metronomev1alpha1.ProviderConfig{}

	if cr.GetProviderConfigReference() == nil {
		return nil, errors.New(errProviderConfigNotSet)
	}

	if err := c.usage.Track(ctx, cr); err != nil {
		return nil, errors.Wrap(err, errFailedToTrackUsage)
	}

	n := types.NamespacedName{Name: cr.GetProviderConfigReference().Name}
	if err := c.client.Get(ctx, n, pc); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	cd := pc.Spec.Credentials
	kc, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.client, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	m := c.newMetronomeClientFn(c.logger, c.baseURL, string(kc))

	return &metronomeExternal{
		logger:    l,
		metronome: m,
	}, nil
}

type metronomeExternal struct {
	logger    logging.Logger
	metronome *metronomeClient.Client
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

	metric, err := e.metronome.GetBillableMetric(id)
	if err != nil {
		// the external name isn't valid
		if errors.Is(err, metronomeClient.ErrBillableMetricInvalidName) {
			return managed.ExternalObservation{ResourceExists: false}, nil
		}
		return managed.ExternalObservation{}, errors.Wrap(err, errFailedToGetBillableMetric)
	}

	if metric == nil {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

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

	res, err := e.metronome.CreateBillableMetric(*req)
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

func (e *metronomeExternal) Delete(_ context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.BillableMetric)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotBillableMetric)
	}

	e.logger.Debug("Deleting")

	id := meta.GetExternalName(cr)
	if id == "" {
		return managed.ExternalDelete{}, nil
	}

	_, err := e.metronome.ArchiveBillableMetric(id)
	if err != nil {
		if errors.Is(err, metronomeClient.ErrBillableMetricAlreadyArchived) {
			return managed.ExternalDelete{}, nil
		}
		return managed.ExternalDelete{}, errors.Wrap(err, errArchiveBillableMetric)
	}

	return managed.ExternalDelete{}, nil
}

func isUpToDate(cr *v1alpha1.BillableMetric, metric *metronomeClient.BillableMetric) (bool, string) {
	spec := &cr.Spec.ForProvider

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
