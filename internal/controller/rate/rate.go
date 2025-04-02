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

package rate

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
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/statemetrics"

	"github.com/redbackthomson/provider-metronome/apis/rate/v1alpha1"
	metronomev1alpha1 "github.com/redbackthomson/provider-metronome/apis/v1alpha1"
	metronomeClient "github.com/redbackthomson/provider-metronome/internal/clients/metronome"
	"github.com/redbackthomson/provider-metronome/internal/converters"
)

const (
	errNotRate              = "managed resource is not a Rate custom resource"
	errProviderConfigNotSet = "provider config is not set"
	errGetProviderConfig    = "cannot get provider config"
	errGetCreds             = "failed to create credentials from provider config"
	errFailedToTrackUsage   = "cannot track provider config usage"
	errGetRate              = "failed to get rate"
	errCreateRate           = "cannot create rate"
	errArchiveRate          = "cannot archive rate"
)

// Setup adds a controller that reconciles Rate managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, baseUrl string) error {
	name := managed.ControllerName(v1alpha1.RateGroupKind)

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
		mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, &v1alpha1.RateList{}, o.MetricOptions.PollStateMetricInterval)); err != nil {
		return err
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.RateGroupVersionKind),
		reconcilerOptions...,
	)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.Rate{}).
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
	cr, ok := mg.(*v1alpha1.Rate)
	if !ok {
		return nil, errors.New(errNotRate)
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

// Observe checks to see if the resource already exists. Metronome doesn't give
// rates a unique ID, so the only way to know if we've already created one is to
// search through all of the existing rates (with a bit of server-side
// filtering) and compare the spec against each of them. If our spec matches an
// existing rate, then we will assume we created it, otherwise we assume it
// doesn't exist.
func (e *metronomeExternal) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Rate)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotRate)
	}

	e.logger.Debug("Observing")

	// no such thing as "deleting" a rate, so don't block on observe
	if cr.DeletionTimestamp != nil {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	var foundRate *metronomeClient.Rate
	nextPage := ""
	for true {
		res, err := e.metronome.GetRates(metronomeClient.GetRatesRequest{
			RateCardID: cr.Spec.ForProvider.RateCardID,
			At:         cr.Spec.ForProvider.StartingAt,
			Selectors: []metronomeClient.RateSelector{{
				PricingGroupValues: cr.Spec.ForProvider.PricingGroupValues,
				ProductID:          cr.Spec.ForProvider.ProductID,
			}},
		}, nextPage)
		if err != nil {
			// the external name isn't valid
			return managed.ExternalObservation{}, errors.Wrap(err, errGetRate)
		}

		if res == nil {
			return managed.ExternalObservation{ResourceExists: false}, nil
		}

		found := false
		for _, r := range res.Data {
			if e.isUpToDate(cr, &r) {
				foundRate = &r
				found = true
			}
		}
		if found {
			break
		}

		nextPage = res.NextPage
		if nextPage == "" {
			break
		}
	}

	if foundRate == nil {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	current := cr.Spec.ForProvider.DeepCopy()
	if err := lateInitialize(&cr.Spec.ForProvider, foundRate); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errGetRate)
	}
	isLateInitialized := !cmp.Equal(current, &cr.Spec.ForProvider)

	converter := &converters.RateConverterImpl{}
	cr.Status.AtProvider = *converter.FromRate(foundRate)
	cr.SetConditions(xpv1.Available())

	return managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        true,
		ResourceLateInitialized: isLateInitialized,
	}, nil
}

func (e *metronomeExternal) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Rate)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotRate)
	}

	e.logger.Debug("Creating")

	converter := &converters.RateConverterImpl{}
	req := converter.FromRateSpec(&cr.Spec.ForProvider)

	_, err := e.metronome.AddRate(*req)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateRate)
	}

	return managed.ExternalCreation{}, nil
}

func (e *metronomeExternal) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	_, ok := mg.(*v1alpha1.Rate)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotRate)
	}

	e.logger.Debug("Updating")

	return managed.ExternalUpdate{}, errors.New("updating a rate is not supported")
}

func (e *metronomeExternal) Delete(_ context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	return managed.ExternalDelete{}, nil
}

func (e *metronomeExternal) isUpToDate(cr *v1alpha1.Rate, r *metronomeClient.Rate) bool {
	spec := &cr.Spec.ForProvider

	converter := &converters.RateConverterImpl{}
	params := converter.FromRateToParameters(r)

	sortTiers := func(a, b v1alpha1.Tier) int {
		if a.Price < b.Price {
			return 1
		} else if a.Price > b.Price {
			return -1
		} else {
			if a.Size < b.Size {
				return 1
			} else if a.Size > b.Size {
				return -1
			}
		}
		return 0
	}

	slices.SortFunc(spec.Tiers, sortTiers)
	slices.SortFunc(params.Tiers, sortTiers)

	// don't compare late initialized fields if they haven't been set
	if spec.CreditTypeID == "" {
		params.CreditTypeID = ""
	}

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmpopts.IgnoreFields(v1alpha1.RateParameters{},
			"RateCardID",
		),
	}

	return cmp.Equal(spec, params, opts...)
}

func lateInitialize(in *v1alpha1.RateParameters, r *metronomeClient.Rate) error {
	if in == nil || r == nil {
		return nil
	}
	in.CreditTypeID = r.Details.CreditType.ID

	return nil
}
