package ratecard

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"

	"github.com/redbackthomson/provider-metronome/apis/ratecard/v1alpha1"
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

type rateCardModifier func(mg *v1alpha1.RateCard)

func rateCard(rm ...rateCardModifier) *v1alpha1.RateCard {
	r := &v1alpha1.RateCard{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testResourceName,
			Namespace: testNamespace,
		},
		Spec: v1alpha1.RateCardSpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{
					Name: providerConfigName,
				},
			},
			ForProvider: v1alpha1.RateCardParameters{},
		},
		Status: v1alpha1.RateCardStatus{},
	}

	meta.SetExternalName(r, "external-name")

	for _, m := range rm {
		m(r)
	}

	return r
}

type notRateCardResource struct {
	resource.Managed
}

type MockRateCardClient struct {
	CreateRateCardFn func(ctx context.Context, reqData metronomeClient.CreateRateCardRequest) (*metronomeClient.CreateRateCardResponse, error)
	GetRateCardFn    func(ctx context.Context, reqData metronomeClient.GetRateCardRequest) (*metronomeClient.GetRateCardResponse, error)
	UpdateRateCardFn func(ctx context.Context, reqData metronomeClient.UpdateRateCardRequest) (*metronomeClient.UpdateRateCardResponse, error)
}

// CreateRateCard implements metronome.RateCardClient.
func (m *MockRateCardClient) CreateRateCard(ctx context.Context, reqData metronomeClient.CreateRateCardRequest) (*metronomeClient.CreateRateCardResponse, error) {
	return m.CreateRateCardFn(ctx, reqData)
}

// GetRateCard implements metronome.RateCardClient.
func (m *MockRateCardClient) GetRateCard(ctx context.Context, reqData metronomeClient.GetRateCardRequest) (*metronomeClient.GetRateCardResponse, error) {
	return m.GetRateCardFn(ctx, reqData)
}

// UpdateRateCard implements metronome.RateCardClient.
func (m *MockRateCardClient) UpdateRateCard(ctx context.Context, reqData metronomeClient.UpdateRateCardRequest) (*metronomeClient.UpdateRateCardResponse, error) {
	return m.UpdateRateCardFn(ctx, reqData)
}

var _ (metronomeClient.RateCardClient) = (*MockRateCardClient)(nil)

func Test_External_Observe(t *testing.T) {
	type args struct {
		metronome metronomeClient.RateCardClient
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
		"NotRateCardResource": {
			args: args{
				mg: notRateCardResource{},
			},
			want: want{
				err: errors.New(errNotRateCard),
			},
		},
		"FailedToGetRateCard": {
			args: args{
				metronome: &MockRateCardClient{
					GetRateCardFn: func(ctx context.Context, reqData metronomeClient.GetRateCardRequest) (*metronomeClient.GetRateCardResponse, error) {
						return nil, errBoom
					},
				},
				mg: rateCard(),
			},
			want: want{
				err: errors.Wrap(errBoom, errGetRateCard),
			},
		},
		"RateCardIsInstantlyDeleted": {
			args: args{
				mg: rateCard(
					func(release *v1alpha1.RateCard) {
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
				metronome: &MockRateCardClient{
					GetRateCardFn: func(ctx context.Context, reqData metronomeClient.GetRateCardRequest) (*metronomeClient.GetRateCardResponse, error) {
						return &metronomeClient.GetRateCardResponse{
							Data: metronomeClient.RateCard{
								ID: "id1", Name: "name",
							},
						}, nil
					},
				},
				mg: rateCard(func(mg *v1alpha1.RateCard) {
					mg.Spec.ForProvider = v1alpha1.RateCardParameters{
						Name: "not-name",
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
				metronome: &MockRateCardClient{
					GetRateCardFn: func(ctx context.Context, reqData metronomeClient.GetRateCardRequest) (*metronomeClient.GetRateCardResponse, error) {
						return &metronomeClient.GetRateCardResponse{
							Data: metronomeClient.RateCard{
								ID: "id1", Name: "name",
							},
						}, nil
					},
				},
				mg: rateCard(func(mg *v1alpha1.RateCard) {
					meta.SetExternalName(mg, "id1")
					mg.Spec.ForProvider = v1alpha1.RateCardParameters{
						Name: "name",
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

			ignoreDiff := cmpopts.IgnoreFields(managed.ExternalObservation{}, "Diff")

			got, gotErr := e.Observe(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("e.Observe(...): -want error, +got error: %s", diff)
			}

			if diff := cmp.Diff(tc.want.out, got, ignoreDiff); diff != "" {
				t.Fatalf("e.Observe(...): -want out, +got out: %s", diff)
			}
		})
	}
}

func Test_External_Create(t *testing.T) {
	type args struct {
		metronome metronomeClient.RateCardClient
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
		"NotRateCardResource": {
			args: args{
				mg: notRateCardResource{},
			},
			want: want{
				err: errors.New(errNotRateCard),
			},
		},
		"FailedToCreateRateCard": {
			args: args{
				metronome: &MockRateCardClient{
					CreateRateCardFn: func(ctx context.Context, reqData metronomeClient.CreateRateCardRequest) (*metronomeClient.CreateRateCardResponse, error) {
						return nil, errBoom
					},
				},
				mg: rateCard(),
			},
			want: want{
				err: errors.Wrap(errBoom, errCreateRateCard),
			},
		},
		"Success": {
			args: args{
				metronome: &MockRateCardClient{
					CreateRateCardFn: func(ctx context.Context, reqData metronomeClient.CreateRateCardRequest) (*metronomeClient.CreateRateCardResponse, error) {
						expected := metronomeClient.CreateRateCardRequest{
							Name: "name",
						}
						if diff := cmp.Diff(expected, reqData); diff != "" {
							t.Errorf("CreateRateCard mismatched: -want req, +got req: %s", diff)
						}

						return &metronomeClient.CreateRateCardResponse{
							Data: metronomeClient.IDOnly{ID: "id"},
						}, nil
					},
				},
				mg: rateCard(func(mg *v1alpha1.RateCard) {
					mg.Spec.ForProvider = v1alpha1.RateCardParameters{
						Name: "name",
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
