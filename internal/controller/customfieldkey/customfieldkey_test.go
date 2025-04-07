package customfieldkey

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"

	"github.com/redbackthomson/provider-metronome/apis/customfieldkey/v1alpha1"
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

type customFieldKeyModifier func(mg *v1alpha1.CustomFieldKey)

func customFieldKey(rm ...customFieldKeyModifier) *v1alpha1.CustomFieldKey {
	r := &v1alpha1.CustomFieldKey{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testResourceName,
			Namespace: testNamespace,
		},
		Spec: v1alpha1.CustomFieldKeySpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{
					Name: providerConfigName,
				},
			},
			ForProvider: v1alpha1.CustomFieldKeyParameters{},
		},
		Status: v1alpha1.CustomFieldKeyStatus{},
	}

	meta.SetExternalName(r, "external-name")

	for _, m := range rm {
		m(r)
	}

	return r
}

type notCustomFieldKeyResource struct {
	resource.Managed
}

type MockCustomFieldKeyClient struct {
	CreateCustomFieldKeyFn func(ctx context.Context, reqData metronomeClient.CreateCustomFieldKeyRequest) error
	DeleteCustomFieldKeyFn func(ctx context.Context, reqData metronomeClient.DeleteCustomFieldKeyRequest) error
	ListCustomFieldKeysFn  func(ctx context.Context, reqData metronomeClient.ListCustomFieldKeysRequest, nextPage string) (*metronomeClient.ListCustomFieldKeysResponse, error)
}

func (m *MockCustomFieldKeyClient) CreateCustomFieldKey(ctx context.Context, reqData metronomeClient.CreateCustomFieldKeyRequest) error {
	return m.CreateCustomFieldKeyFn(ctx, reqData)
}

func (m *MockCustomFieldKeyClient) DeleteCustomFieldKey(ctx context.Context, reqData metronomeClient.DeleteCustomFieldKeyRequest) error {
	return m.DeleteCustomFieldKeyFn(ctx, reqData)
}

func (m *MockCustomFieldKeyClient) ListCustomFieldKeys(ctx context.Context, reqData metronomeClient.ListCustomFieldKeysRequest, nextPage string) (*metronomeClient.ListCustomFieldKeysResponse, error) {
	return m.ListCustomFieldKeysFn(ctx, reqData, nextPage)
}

var _ (metronomeClient.CustomFieldKeyClient) = (*MockCustomFieldKeyClient)(nil)

func Test_External_Observe(t *testing.T) {
	type args struct {
		metronome metronomeClient.CustomFieldKeyClient
		mg        resource.Managed
	}
	type want struct {
		out managed.ExternalObservation
		err error
	}
	cases := map[string]struct {
		args
		want
	}{
		"NotCustomFieldKeyResource": {
			args: args{
				mg: notCustomFieldKeyResource{},
			},
			want: want{
				err: errors.New(errNotCustomFieldKey),
			},
		},
		"NoCustomFieldKeyExists": {
			args: args{
				metronome: &MockCustomFieldKeyClient{
					ListCustomFieldKeysFn: func(ctx context.Context, reqData metronomeClient.ListCustomFieldKeysRequest, nextPage string) (*metronomeClient.ListCustomFieldKeysResponse, error) {
						return &metronomeClient.ListCustomFieldKeysResponse{}, nil
					},
				},
				mg: customFieldKey(),
			},
			want: want{
				out: managed.ExternalObservation{ResourceExists: false},
				err: nil,
			},
		},
		"FailedToGetCustomFieldKey": {
			args: args{
				metronome: &MockCustomFieldKeyClient{
					ListCustomFieldKeysFn: func(ctx context.Context, reqData metronomeClient.ListCustomFieldKeysRequest, nextPage string) (*metronomeClient.ListCustomFieldKeysResponse, error) {
						return nil, errBoom
					},
				},
				mg: customFieldKey(),
			},
			want: want{
				err: errors.Wrap(errBoom, errGetCustomFieldKey),
			},
		},
		"CustomFieldKeyIsInstantlyDeleted": {
			args: args{
				mg: customFieldKey(
					func(release *v1alpha1.CustomFieldKey) {
						now := metav1.Now()
						release.SetDeletionTimestamp(&now)
					},
				),
			},
			want: want{
				out: managed.ExternalObservation{ResourceExists: false},
			},
		},
		"NotUpToDate": {
			args: args{
				metronome: &MockCustomFieldKeyClient{
					ListCustomFieldKeysFn: func(ctx context.Context, reqData metronomeClient.ListCustomFieldKeysRequest, nextPage string) (*metronomeClient.ListCustomFieldKeysResponse, error) {
						return &metronomeClient.ListCustomFieldKeysResponse{
							Data: []metronomeClient.CustomFieldKey{
								{Key: "key1", Entity: "entity1", EnforceUniqueness: false},
							},
						}, nil
					},
				},
				mg: customFieldKey(func(mg *v1alpha1.CustomFieldKey) {
					mg.Spec.ForProvider = v1alpha1.CustomFieldKeyParameters{
						Key:               "key1",
						Entity:            "entity1",
						EnforceUniqueness: true,
					}
				}),
			},
			want: want{
				out: managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: false},
				err: nil,
			},
		},
		"UpToDate": {
			args: args{
				metronome: &MockCustomFieldKeyClient{
					ListCustomFieldKeysFn: func(ctx context.Context, reqData metronomeClient.ListCustomFieldKeysRequest, nextPage string) (*metronomeClient.ListCustomFieldKeysResponse, error) {
						switch nextPage {
						case "":
							return &metronomeClient.ListCustomFieldKeysResponse{
								Data: []metronomeClient.CustomFieldKey{
									{Key: "key0", Entity: "entity0"},
								},
								NextPage: "1",
							}, nil
						case "1":
							return &metronomeClient.ListCustomFieldKeysResponse{
								Data: []metronomeClient.CustomFieldKey{
									{Key: "key1", Entity: "entity1", EnforceUniqueness: true},
								},
							}, nil
						default:
							return nil, errBoom
						}
					},
				},
				mg: customFieldKey(func(mg *v1alpha1.CustomFieldKey) {
					mg.Spec.ForProvider = v1alpha1.CustomFieldKeyParameters{
						Key:               "key1",
						Entity:            "entity1",
						EnforceUniqueness: true,
					}
				}),
			},
			want: want{
				out: managed.ExternalObservation{ResourceExists: true, ResourceUpToDate: true},
				err: nil,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &metronomeExternal{
				logger:    logging.NewNopLogger(),
				metronome: tc.args.metronome,
			}
			got, gotErr := e.Observe(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("e.Observe(...): -want error, +got error: %s", diff)
			}

			if diff := cmp.Diff(tc.want.out, got); diff != "" {
				t.Fatalf("e.Observe(...): -want out, +got out: %s", diff)
			}
		})
	}
}

func Test_External_Create(t *testing.T) {
	type args struct {
		metronome metronomeClient.CustomFieldKeyClient
		mg        resource.Managed
	}
	type want struct {
		out managed.ExternalCreation
		err error
	}
	cases := map[string]struct {
		args
		want
	}{
		"NotCustomFieldKeyResource": {
			args: args{
				mg: notCustomFieldKeyResource{},
			},
			want: want{
				err: errors.New(errNotCustomFieldKey),
			},
		},
		"FailedToCreateCustomFieldKey": {
			args: args{
				metronome: &MockCustomFieldKeyClient{
					CreateCustomFieldKeyFn: func(ctx context.Context, reqData metronomeClient.CreateCustomFieldKeyRequest) error {
						return errBoom
					},
				},
				mg: customFieldKey(),
			},
			want: want{
				err: errors.Wrap(errBoom, errCreateCustomFieldKey),
			},
		},
		"Success": {
			args: args{
				metronome: &MockCustomFieldKeyClient{
					CreateCustomFieldKeyFn: func(ctx context.Context, reqData metronomeClient.CreateCustomFieldKeyRequest) error {
						expected := metronomeClient.CreateCustomFieldKeyRequest{
							Key:               "key",
							Entity:            "entity",
							EnforceUniqueness: true,
						}
						if diff := cmp.Diff(expected, reqData); diff != "" {
							t.Errorf("CreateCustomFieldKey mismatched: -want req, +got req: %s", diff)
						}

						return nil
					},
				},
				mg: customFieldKey(func(mg *v1alpha1.CustomFieldKey) {
					mg.Spec.ForProvider = v1alpha1.CustomFieldKeyParameters{
						Key:               "key",
						Entity:            "entity",
						EnforceUniqueness: true,
					}
				}),
			},
			want: want{
				out: managed.ExternalCreation{},
				err: nil,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &metronomeExternal{
				logger:    logging.NewNopLogger(),
				metronome: tc.args.metronome,
			}
			got, gotErr := e.Create(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("e.Create(...): -want error, +got error: %s", diff)
			}

			if diff := cmp.Diff(tc.want.out, got); diff != "" {
				t.Fatalf("e.Create(...): -want out, +got out: %s", diff)
			}
		})
	}
}

func Test_External_Delete(t *testing.T) {
	type args struct {
		metronome metronomeClient.CustomFieldKeyClient
		mg        resource.Managed
	}
	type want struct {
		out managed.ExternalDelete
		err error
	}
	cases := map[string]struct {
		args
		want
	}{
		"NotCustomFieldKeyResource": {
			args: args{
				mg: notCustomFieldKeyResource{},
			},
			want: want{
				err: errors.New(errNotCustomFieldKey),
			},
		},
		"FailedToDeleteCustomFieldKey": {
			args: args{
				metronome: &MockCustomFieldKeyClient{
					DeleteCustomFieldKeyFn: func(ctx context.Context, reqData metronomeClient.DeleteCustomFieldKeyRequest) error {
						return errBoom
					},
				},
				mg: customFieldKey(),
			},
			want: want{
				err: errors.Wrap(errBoom, errDeleteCustomFieldKey),
			},
		},
		"Success": {
			args: args{
				metronome: &MockCustomFieldKeyClient{
					DeleteCustomFieldKeyFn: func(ctx context.Context, reqData metronomeClient.DeleteCustomFieldKeyRequest) error {
						expected := metronomeClient.DeleteCustomFieldKeyRequest{
							Key:    "key",
							Entity: "entity",
						}
						if diff := cmp.Diff(expected, reqData); diff != "" {
							t.Errorf("DeleteCustomFieldKey mismatched: -want req, +got req: %s", diff)
						}

						return nil
					},
				},
				mg: customFieldKey(func(mg *v1alpha1.CustomFieldKey) {
					mg.Spec.ForProvider = v1alpha1.CustomFieldKeyParameters{
						Key:               "key",
						Entity:            "entity",
						EnforceUniqueness: true,
					}
				}),
			},
			want: want{
				out: managed.ExternalDelete{},
				err: nil,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &metronomeExternal{
				logger:    logging.NewNopLogger(),
				metronome: tc.args.metronome,
			}
			got, gotErr := e.Delete(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("e.Delete(...): -want error, +got error: %s", diff)
			}

			if diff := cmp.Diff(tc.want.out, got); diff != "" {
				t.Fatalf("e.Delete(...): -want out, +got out: %s", diff)
			}
		})
	}
}
