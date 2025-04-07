package rate

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"

	"github.com/redbackthomson/provider-metronome/apis/rate/v1alpha1"
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

type rateModifier func(release *v1alpha1.Rate)

func rate(rm ...rateModifier) *v1alpha1.Rate {
	r := &v1alpha1.Rate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testResourceName,
			Namespace: testNamespace,
		},
		Spec: v1alpha1.RateSpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{
					Name: providerConfigName,
				},
			},
			ForProvider: v1alpha1.RateParameters{},
		},
		Status: v1alpha1.RateStatus{},
	}

	for _, m := range rm {
		m(r)
	}

	return r
}

func fullyPopulate(release *v1alpha1.Rate) {
	release.Spec.ForProvider = v1alpha1.RateParameters{
		RateCardID: "rate-card-id",
		ProductID:  "product-id",
		StartingAt: "starting-at",
		Entitled:   true,
		RateType:   "rate-type",
		Price:      1.01,
		PricingGroupValues: map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
		CommitRate: &v1alpha1.CommitRate{
			RateType: "commit-rate-type",
			Price:    1.02,
			Tiers: []v1alpha1.Tier{{
				Price: 1.03,
				Size:  10,
			}, {
				Price: 1.04,
				Size:  20,
			}},
		},
		CreditTypeID: "credit-type-id",
		EndingBefore: "ending-before",
		IsProrated:   true,
		Quantity:     1.05,
		Tiers: []v1alpha1.Tier{{
			Price: 1.06,
			Size:  30,
		}, {
			Price: 1.07,
			Size:  40,
		}},
		UseListPrices: true,
	}
}

var fullyPopulated = metronomeClient.Rate{
	Entitled:     true,
	ProductID:    "product-id",
	ProductName:  "product-name",
	StartingAt:   "starting-at",
	EndingBefore: "ending-before",
	PricingGroupValues: map[string]string{
		"key1": "val1",
		"key2": "val2",
	},
	ProductTags: []string{
		"tag1", "tag2",
	},
	ProductCustomField: map[string]string{
		"custom1": "val1",
		"custom2": "val2",
	},
	CommitRate: &metronomeClient.CommitRate{
		RateType: "commit-rate-type",
		Price:    1.02,
		Tiers: []metronomeClient.Tier{{
			Price: 1.03,
			Size:  10,
		}, {
			Price: 1.04,
			Size:  20,
		}},
	},
	Details: metronomeClient.RateDetails{
		RateType:   "rate-type",
		IsProrated: true,
		Price:      1.01,
		PricingGroupValues: map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
		Quantity: 1.05,
		Tiers: []metronomeClient.Tier{{
			Price: 1.06,
			Size:  30,
		}, {
			Price: 1.07,
			Size:  40,
		}},
		UseListPrices: true,
		CreditType: metronomeClient.CreditType{
			ID:   "credit-type-id",
			Name: "credit-type-name",
		},
	},
}

type notRateResource struct {
	resource.Managed
}

type MockRateClient struct {
	GetRatesFn func(ctx context.Context, reqData metronomeClient.GetRatesRequest, nextPage string) (*metronomeClient.GetRatesResponse, error)
	AddRateFn  func(ctx context.Context, reqData metronomeClient.AddRateRequest) (*metronomeClient.AddRateResponse, error)
}

func (m *MockRateClient) AddRate(ctx context.Context, reqData metronomeClient.AddRateRequest) (*metronomeClient.AddRateResponse, error) {
	return m.AddRateFn(ctx, reqData)
}

func (m *MockRateClient) GetRates(ctx context.Context, reqData metronomeClient.GetRatesRequest, nextPage string) (*metronomeClient.GetRatesResponse, error) {
	return m.GetRatesFn(ctx, reqData, nextPage)
}

var _ (metronomeClient.RateClient) = (*MockRateClient)(nil)

func Test_External_Observe(t *testing.T) {
	type args struct {
		metronome metronomeClient.RateClient
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
		"NotRateResource": {
			args: args{
				mg: notRateResource{},
			},
			want: want{
				err: errors.New(errNotRate),
			},
		},
		"NoRateExists": {
			args: args{
				metronome: &MockRateClient{
					GetRatesFn: func(ctx context.Context, reqData metronomeClient.GetRatesRequest, nextPage string) (*metronomeClient.GetRatesResponse, error) {
						return &metronomeClient.GetRatesResponse{Data: []metronomeClient.Rate{}}, nil
					},
				},
				mg: rate(),
			},
			want: want{
				out: managed.ExternalObservation{ResourceExists: false},
				err: nil,
			},
		},
		"FailedToGetRate": {
			args: args{
				metronome: &MockRateClient{
					GetRatesFn: func(ctx context.Context, reqData metronomeClient.GetRatesRequest, nextPage string) (*metronomeClient.GetRatesResponse, error) {
						return nil, errBoom
					},
				},
				mg: rate(),
			},
			want: want{
				err: errors.Wrap(errBoom, errGetRate),
			},
		},
		"RateIsInstantlyDeleted": {
			args: args{
				metronome: &MockRateClient{
					GetRatesFn: func(ctx context.Context, reqData metronomeClient.GetRatesRequest, nextPage string) (*metronomeClient.GetRatesResponse, error) {
						return nil, errBoom
					},
				},
				mg: rate(
					func(release *v1alpha1.Rate) {
						now := metav1.Now()
						release.SetDeletionTimestamp(&now)
					},
				),
			},
			want: want{
				out: managed.ExternalObservation{ResourceExists: false},
			},
		},
		"UpToDate": {
			args: args{
				metronome: &MockRateClient{
					GetRatesFn: func(ctx context.Context, reqData metronomeClient.GetRatesRequest, nextPage string) (*metronomeClient.GetRatesResponse, error) {
						switch nextPage {
						case "":
							return &metronomeClient.GetRatesResponse{
								Data: []metronomeClient.Rate{{
									ProductID: "first",
								}},
								NextPage: "1",
							}, nil
						case "1":
							return &metronomeClient.GetRatesResponse{
								Data:     []metronomeClient.Rate{fullyPopulated},
								NextPage: "1",
							}, nil
						default:
							return nil, nil
						}
					},
				},
				mg: rate(fullyPopulate),
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
		metronome metronomeClient.RateClient
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
		"NotRateResource": {
			args: args{
				mg: notRateResource{},
			},
			want: want{
				err: errors.New(errNotRate),
			},
		},
		"FailedToCreateRate": {
			args: args{
				metronome: &MockRateClient{
					AddRateFn: func(ctx context.Context, reqData metronomeClient.AddRateRequest) (*metronomeClient.AddRateResponse, error) {
						return nil, errBoom
					},
				},
				mg: rate(),
			},
			want: want{
				err: errors.Wrap(errBoom, errCreateRate),
			},
		},
		"Success": {
			args: args{
				metronome: &MockRateClient{
					AddRateFn: func(ctx context.Context, reqData metronomeClient.AddRateRequest) (*metronomeClient.AddRateResponse, error) {
						expected := metronomeClient.AddRateRequest{
							Entitled:      true,
							ProductID:     "product-id",
							RateCardID:    "rate-card-id",
							RateType:      "rate-type",
							StartingAt:    "starting-at",
							CreditTypeID:  "credit-type-id",
							EndingBefore:  "ending-before",
							IsProrated:    true,
							Price:         1.01,
							Quantity:      1.05,
							UseListPrices: true,
							Tiers: []metronomeClient.Tier{
								{Price: 1.06, Size: 30},
								{Price: 1.07, Size: 40},
							},
							PricingGroupValues: map[string]string{
								"key1": "val1",
								"key2": "val2",
							},
							CommitRate: &metronomeClient.CommitRate{
								RateType: "commit-rate-type",
								Price:    1.02,
								Tiers: []metronomeClient.Tier{
									{Price: 1.03, Size: 10},
									{Price: 1.04, Size: 20},
								},
							},
						}
						if diff := cmp.Diff(expected, reqData); diff != "" {
							t.Errorf("AddRateRequest mismatched: -want req, +got req: %s", diff)
						}

						return &metronomeClient.AddRateResponse{}, nil
					},
				},
				mg: rate(fullyPopulate),
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
