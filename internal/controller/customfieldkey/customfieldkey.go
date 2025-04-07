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

package customfieldkey

import (
	"context"

	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/feature"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/statemetrics"

	"github.com/redbackthomson/provider-metronome/apis/customfieldkey/v1alpha1"
	metronomev1alpha1 "github.com/redbackthomson/provider-metronome/apis/v1alpha1"
	metronomeClient "github.com/redbackthomson/provider-metronome/internal/clients/metronome"
	"github.com/redbackthomson/provider-metronome/internal/connector"
	"github.com/redbackthomson/provider-metronome/internal/converters"
)

const (
	errNotCustomFieldKey    = "managed resource is not a CustomFieldKey custom resource"
	errGetCustomFieldKey    = "failed to get custom field key"
	errCreateCustomFieldKey = "failed to create custom field key"
	errDeleteCustomFieldKey = "failed to delete custom field key"
)

// Setup adds a controller that reconciles CustomFieldKey managed resources.
func Setup(mgr ctrl.Manager, o controller.Options, baseUrl string) error {
	name := managed.ControllerName(v1alpha1.CustomFieldKeyGroupKind)

	reconcilerOptions := []managed.ReconcilerOption{
		managed.WithExternalConnecter(
			&connector.Connector[*v1alpha1.CustomFieldKey, *metronomeExternal]{
				Logger:               o.Logger,
				Client:               mgr.GetClient(),
				Usage:                resource.NewProviderConfigUsageTracker(mgr.GetClient(), &metronomev1alpha1.ProviderConfigUsage{}),
				BaseURL:              baseUrl,
				NewMetronomeClientFn: metronomeClient.New,
				NewExternalClientFn: func(log logging.Logger, client *metronomeClient.Client) *metronomeExternal {
					return &metronomeExternal{
						logger:    o.Logger,
						metronome: client.CustomFieldKey(),
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
		mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, &v1alpha1.CustomFieldKeyList{}, o.MetricOptions.PollStateMetricInterval)); err != nil {
		return err
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.CustomFieldKeyGroupVersionKind),
		reconcilerOptions...,
	)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.CustomFieldKey{}).
		WithOptions(o.ForControllerRuntime()).
		Complete(r)
}

type metronomeExternal struct {
	logger    logging.Logger
	metronome metronomeClient.CustomFieldKeyClient
}

func (e *metronomeExternal) Disconnect(ctx context.Context) error {
	return nil
}

// Observe checks to see if the resource already exists. Metronome doesn't give
// custom field keys a unique ID, so the only way to know if we've already
// created one is to search through all of the existing custom fiel dkeys (with
// a bit of server-side filtering) and compare the spec against each of them. If
// our spec matches an existing custom field key, then we will assume we created
// it, otherwise we assume it doesn't exist.
func (e *metronomeExternal) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.CustomFieldKey)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotCustomFieldKey)
	}

	e.logger.Debug("Observing")

	// no such thing as "deleting" a customfieldkey, so don't block on observe
	if cr.DeletionTimestamp != nil {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	var foundCustomFieldKey *metronomeClient.CustomFieldKey
	nextPage := ""
	for {
		res, err := e.metronome.ListCustomFieldKeys(ctx, metronomeClient.ListCustomFieldKeysRequest{
			Entities: []string{cr.Spec.ForProvider.Entity},
		}, nextPage)
		if err != nil {
			return managed.ExternalObservation{}, errors.Wrap(err, errGetCustomFieldKey)
		}

		if res == nil {
			return managed.ExternalObservation{ResourceExists: false}, nil
		}

		found := false
		for _, r := range res.Data {
			if r.Key == cr.Spec.ForProvider.Key && r.Entity == cr.Spec.ForProvider.Entity {
				foundCustomFieldKey = &r
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

	if foundCustomFieldKey == nil {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}

	converter := &converters.CustomFieldKeyConverterImpl{}
	cr.Status.AtProvider = *converter.FromCustomFieldKey(foundCustomFieldKey)
	cr.SetConditions(xpv1.Available())

	upToDate := foundCustomFieldKey.EnforceUniqueness == cr.Spec.ForProvider.EnforceUniqueness

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: upToDate,
	}, nil
}

func (e *metronomeExternal) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.CustomFieldKey)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotCustomFieldKey)
	}

	e.logger.Debug("Creating")

	converter := &converters.CustomFieldKeyConverterImpl{}
	req := converter.FromCustomFieldKeySpec(&cr.Spec.ForProvider)

	err := e.metronome.CreateCustomFieldKey(ctx, *req)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateCustomFieldKey)
	}

	return managed.ExternalCreation{}, nil
}

func (e *metronomeExternal) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	_, ok := mg.(*v1alpha1.CustomFieldKey)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotCustomFieldKey)
	}

	e.logger.Debug("Updating")

	return managed.ExternalUpdate{}, errors.New("updating a custom field key is not supported")
}

func (e *metronomeExternal) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.CustomFieldKey)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotCustomFieldKey)
	}

	e.logger.Debug("Deleting")

	err := e.metronome.DeleteCustomFieldKey(ctx, metronomeClient.DeleteCustomFieldKeyRequest{
		Entity: cr.Spec.ForProvider.Entity,
		Key:    cr.Spec.ForProvider.Key,
	})
	if err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, errDeleteCustomFieldKey)
	}

	return managed.ExternalDelete{}, nil
}
