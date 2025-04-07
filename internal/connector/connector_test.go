package connector

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"

	"github.com/redbackthomson/provider-metronome/apis/billablemetric/v1alpha1"
	metronomev1alpha1 "github.com/redbackthomson/provider-metronome/apis/v1alpha1"
	metronomeClient "github.com/redbackthomson/provider-metronome/internal/clients/metronome"
)

const (
	providerConfigName = "metronome-test"
	testResourceName   = "test-resource"
	testNamespace      = "testns"
)

var (
	errBoom = errors.New("boom")
)

// use billable metric for testing expected resource
type expectedResource = v1alpha1.BillableMetric

func billableMetric() *v1alpha1.BillableMetric {
	return &v1alpha1.BillableMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testResourceName,
			Namespace: testNamespace,
		},
		Spec: v1alpha1.BillableMetricSpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{
					Name: providerConfigName,
				},
			},
			ForProvider: v1alpha1.BillableMetricParameters{},
		},
		Status: v1alpha1.BillableMetricStatus{},
	}
}

type unexpectedResource struct {
	resource.Managed
}

type mockExternalClient struct {
	managed.ExternalClient
}

var _ (managed.ExternalClient) = (*mockExternalClient)(nil)

func (c *mockExternalClient) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	return managed.ExternalObservation{}, nil
}

func (c *mockExternalClient) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	return managed.ExternalCreation{}, nil
}

func (c *mockExternalClient) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil
}

func (c *mockExternalClient) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	return managed.ExternalDelete{}, nil
}

func (c *mockExternalClient) Disconnect(ctx context.Context) error {
	return nil
}

func Test_Connector_Connect(t *testing.T) {
	providerConfig := metronomev1alpha1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: providerConfigName},
		Spec: metronomev1alpha1.ProviderConfigSpec{
			Credentials: metronomev1alpha1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						SecretReference: xpv1.SecretReference{
							Name:      "creds",
							Namespace: testNamespace,
						},
						Key: "auth",
					},
				},
			},
		},
	}

	type args struct {
		client client.Client
		usage  resource.Tracker
		mg     resource.Managed

		baseURL              string
		newMetronomeClientFn func(log logging.Logger, baseURL, authToken string) (*metronomeClient.Client, error)
		newExternalClientFn  func(log logging.Logger, client *metronomeClient.Client) *mockExternalClient
	}
	type want struct {
		err error
	}
	cases := map[string]struct {
		args
		want
	}{
		"UnexpectedResource": {
			args: args{
				mg: unexpectedResource{},
			},
			want: want{
				err: errors.New(errResourceWrongType),
			},
		},
		"FailedToTrackUsage": {
			args: args{
				usage: resource.TrackerFn(func(ctx context.Context, mg resource.Managed) error { return errBoom }),
				mg:    billableMetric(),
			},
			want: want{
				err: errors.Wrap(errBoom, errFailedToTrackUsage),
			},
		},
		"FailedToGetProvider": {
			args: args{
				client: &test.MockClient{
					MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
						if key.Name == providerConfigName {
							*obj.(*metronomev1alpha1.ProviderConfig) = providerConfig
							return errBoom
						}
						return nil
					},
				},
				usage: resource.TrackerFn(func(ctx context.Context, mg resource.Managed) error { return nil }),
				mg:    billableMetric(),
			},
			want: want{
				err: errors.Wrap(errBoom, errGetProviderConfig),
			},
		},
		"FailedToCreateNewMetronomeClient": {
			args: args{
				client: &test.MockClient{
					MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
						switch t := obj.(type) {
						case *metronomev1alpha1.ProviderConfig:
							*t = providerConfig
						case *corev1.Secret:
							*t = corev1.Secret{
								Data: map[string][]byte{
									"auth": []byte("def456"),
								},
							}
						default:
							return errBoom
						}
						return nil
					},
					MockStatusUpdate: func(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
						return nil
					},
				},
				newMetronomeClientFn: func(log logging.Logger, baseURL, authToken string) (*metronomeClient.Client, error) {
					return nil, errBoom
				},
				usage: resource.TrackerFn(func(ctx context.Context, mg resource.Managed) error { return nil }),
				mg:    billableMetric(),
			},
			want: want{
				err: errors.Wrap(errBoom, errConnectToMetronome),
			},
		},
		"Success": {
			args: args{
				client: &test.MockClient{
					MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
						switch t := obj.(type) {
						case *metronomev1alpha1.ProviderConfig:
							*t = providerConfig
						case *corev1.Secret:
							*t = corev1.Secret{
								Data: map[string][]byte{
									"auth": []byte("def456"),
								},
							}
						default:
							return errBoom
						}
						return nil
					},
					MockStatusUpdate: func(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
						return nil
					},
				},
				baseURL: "abc123",
				newMetronomeClientFn: func(log logging.Logger, baseURL, authToken string) (*metronomeClient.Client, error) {
					if baseURL != "abc123" {
						t.Errorf("unexpected base URL: %s", baseURL)
					}
					if authToken != "def456" {
						t.Errorf("unexpected auth token: %s", authToken)
					}
					return &metronomeClient.Client{}, nil
				},
				newExternalClientFn: func(log logging.Logger, client *metronomeClient.Client) *mockExternalClient {
					return &mockExternalClient{}
				},
				usage: resource.TrackerFn(func(ctx context.Context, mg resource.Managed) error { return nil }),
				mg:    billableMetric(),
			},
			want: want{
				err: nil,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c := &Connector[*expectedResource, *mockExternalClient]{
				Logger:  logging.NewNopLogger(),
				Client:  tc.client,
				Usage:   tc.usage,
				BaseURL: tc.baseURL,

				NewMetronomeClientFn: tc.newMetronomeClientFn,
				NewExternalClientFn:  tc.newExternalClientFn,
			}
			_, gotErr := c.Connect(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("Connect(...): -want error, +got error: %s", diff)
			}
		})
	}
}
