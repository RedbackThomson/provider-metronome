package connector

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metronomev1alpha1 "github.com/redbackthomson/provider-metronome/apis/v1alpha1"
	metronomeClient "github.com/redbackthomson/provider-metronome/internal/clients/metronome"
)

const (
	errResourceWrongType    = "resource does not match connector type"
	errProviderConfigNotSet = "provider config is not set"
	errGetProviderConfig    = "cannot get provider config"
	errGetCreds             = "failed to create credentials from provider config"
	errFailedToTrackUsage   = "cannot track provider config usage"
	errConnectToMetronome   = "error connecting to Metronome"
)

type Connector[R resource.Managed, T managed.ExternalClient] struct {
	BaseURL string
	Logger  logging.Logger
	Client  client.Client
	Usage   resource.Tracker

	NewMetronomeClientFn func(log logging.Logger, baseURL, authToken string) (*metronomeClient.Client, error)
	NewExternalClientFn  func(log logging.Logger, client *metronomeClient.Client) T
}

func (c *Connector[R, T]) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) { //nolint:gocyclo
	cr, ok := mg.(R)
	if !ok {
		return nil, errors.New(errResourceWrongType)
	}
	l := c.Logger.WithValues("request", cr.GetName())

	l.Debug("Connecting")

	pc := &metronomev1alpha1.ProviderConfig{}

	if cr.GetProviderConfigReference() == nil {
		return nil, errors.New(errProviderConfigNotSet)
	}

	if err := c.Usage.Track(ctx, cr); err != nil {
		return nil, errors.Wrap(err, errFailedToTrackUsage)
	}

	n := types.NamespacedName{Name: cr.GetProviderConfigReference().Name}
	if err := c.Client.Get(ctx, n, pc); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	cd := pc.Spec.Credentials
	kc, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.Client, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	m, err := c.NewMetronomeClientFn(c.Logger, c.BaseURL, string(kc))
	if err != nil {
		return nil, errors.Wrap(err, errConnectToMetronome)
	}

	return c.NewExternalClientFn(c.Logger, m), nil
}
