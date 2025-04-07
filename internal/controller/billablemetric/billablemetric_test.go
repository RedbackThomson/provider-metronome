package billablemetric

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"

	"github.com/redbackthomson/provider-metronome/apis/billablemetric/v1alpha1"
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

type billableMetricModifier func(mg *v1alpha1.BillableMetric)

func billableMetric(rm ...billableMetricModifier) *v1alpha1.BillableMetric {
	r := &v1alpha1.BillableMetric{
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

	meta.SetExternalName(r, "external-name")

	for _, m := range rm {
		m(r)
	}

	return r
}

type notBillableMetricResource struct {
	resource.Managed
}

type MockBillableMetricClient struct {
	ArchiveBillableMetricFn func(ctx context.Context, id string) (*metronomeClient.ArchiveBillableMetricResponse, error)
	CreateBillableMetricFn  func(ctx context.Context, reqData metronomeClient.CreateBillableMetricRequest) (*metronomeClient.CreateBillableMetricResponse, error)
	GetBillableMetricFn     func(ctx context.Context, id string) (*metronomeClient.GetBillableMetricResponse, error)
	ListBillableMetricsFn   func(ctx context.Context) (*metronomeClient.ListBillableMetricsResponse, error)
	UpdateBillableMetricFn  func(ctx context.Context, id string, reqData metronomeClient.UpdateBillableMetricRequest) (*metronomeClient.UpdateBillableMetricResponse, error)
}

func (m *MockBillableMetricClient) ArchiveBillableMetric(ctx context.Context, id string) (*metronomeClient.ArchiveBillableMetricResponse, error) {
	return m.ArchiveBillableMetricFn(ctx, id)
}

func (m *MockBillableMetricClient) CreateBillableMetric(ctx context.Context, reqData metronomeClient.CreateBillableMetricRequest) (*metronomeClient.CreateBillableMetricResponse, error) {
	return m.CreateBillableMetricFn(ctx, reqData)
}

func (m *MockBillableMetricClient) GetBillableMetric(ctx context.Context, id string) (*metronomeClient.GetBillableMetricResponse, error) {
	return m.GetBillableMetricFn(ctx, id)
}

func (m *MockBillableMetricClient) ListBillableMetrics(ctx context.Context) (*metronomeClient.ListBillableMetricsResponse, error) {
	return m.ListBillableMetricsFn(ctx)
}

func (m *MockBillableMetricClient) UpdateBillableMetric(ctx context.Context, id string, reqData metronomeClient.UpdateBillableMetricRequest) (*metronomeClient.UpdateBillableMetricResponse, error) {
	return m.UpdateBillableMetricFn(ctx, id, reqData)
}

var _ (metronomeClient.BillableMetricClient) = (*MockBillableMetricClient)(nil)

func Test_External_Observe(t *testing.T) {
	type args struct {
		metronome metronomeClient.BillableMetricClient
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
		"NotBillableMetricResource": {
			args: args{
				mg: notBillableMetricResource{},
			},
			want: want{
				err: errors.New(errNotBillableMetric),
			},
		},
		"NoExternalName": {
			args: args{
				metronome: &MockBillableMetricClient{
					GetBillableMetricFn: func(ctx context.Context, id string) (*metronomeClient.GetBillableMetricResponse, error) {
						return nil, nil
					},
				},
				mg: billableMetric(),
			},
			want: want{
				out: managed.ExternalObservation{ResourceExists: false},
				err: nil,
			},
		},
		"NoBillableMetricExists": {
			args: args{
				metronome: &MockBillableMetricClient{
					GetBillableMetricFn: func(ctx context.Context, id string) (*metronomeClient.GetBillableMetricResponse, error) {
						return nil, nil
					},
				},
				mg: billableMetric(func(mg *v1alpha1.BillableMetric) {
					meta.SetExternalName(mg, "")
				}),
			},
			want: want{
				out: managed.ExternalObservation{ResourceExists: false},
				err: nil,
			},
		},
		"InvalidName": {
			args: args{
				metronome: &MockBillableMetricClient{
					GetBillableMetricFn: func(ctx context.Context, id string) (*metronomeClient.GetBillableMetricResponse, error) {
						return nil, metronomeClient.ErrBillableMetricInvalidName
					},
				},
				mg: billableMetric(),
			},
			want: want{
				out: managed.ExternalObservation{ResourceExists: false},
			},
		},
		"FailedToGetBillableMetric": {
			args: args{
				metronome: &MockBillableMetricClient{
					GetBillableMetricFn: func(ctx context.Context, id string) (*metronomeClient.GetBillableMetricResponse, error) {
						return nil, errBoom
					},
				},
				mg: billableMetric(),
			},
			want: want{
				err: errors.Wrap(errBoom, errGetBillableMetric),
			},
		},
		"UpToDate": {
			args: args{
				metronome: &MockBillableMetricClient{
					GetBillableMetricFn: func(ctx context.Context, id string) (*metronomeClient.GetBillableMetricResponse, error) {
						return &metronomeClient.GetBillableMetricResponse{
							Data: metronomeClient.BillableMetric{
								ID:              "id",
								Name:            "name",
								AggregationType: metronomeClient.AggregationMax,
								AggregationKey:  "agg-key",
								EventTypeFilter: metronomeClient.EventTypeFilter{
									InValues: []string{"in", "in2"},
								},
								PropertyFilters: []metronomeClient.PropertyFilter{{
									Name:     "prop-filter-name",
									Exists:   ptr.To(true),
									InValues: []string{"prop-filter-in-1", "prop-filter-in-2"},
								}},
								GroupKeys: [][]string{{"group-key-1", "group-key-2"}, {"group-2-key-1"}},
							},
						}, nil
					},
				},
				mg: billableMetric(func(mg *v1alpha1.BillableMetric) {
					mg.Spec.ForProvider = v1alpha1.BillableMetricParameters{
						Name:            "name",
						AggregationType: v1alpha1.AggregationTypeMax,
						AggregationKey:  "agg-key",
						EventTypeFilter: v1alpha1.EventTypeFilter{
							InValues: []string{"in", "in2"},
						},
						PropertyFilters: []v1alpha1.PropertyFilter{{
							Name:     "prop-filter-name",
							Exists:   ptr.To(true),
							InValues: []string{"prop-filter-in-1", "prop-filter-in-2"},
						}},
						GroupKeys: [][]string{{"group-key-1", "group-key-2"}, {"group-2-key-1"}},
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
		metronome metronomeClient.BillableMetricClient
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
		"NotBillableMetricResource": {
			args: args{
				mg: notBillableMetricResource{},
			},
			want: want{
				err: errors.New(errNotBillableMetric),
			},
		},
		"FailedToCreateBillableMetric": {
			args: args{
				metronome: &MockBillableMetricClient{
					CreateBillableMetricFn: func(ctx context.Context, reqData metronomeClient.CreateBillableMetricRequest) (*metronomeClient.CreateBillableMetricResponse, error) {
						return nil, errBoom
					},
				},
				mg: billableMetric(),
			},
			want: want{
				err: errors.Wrap(errBoom, errCreateBillableMetric),
			},
		},
		"Success": {
			args: args{
				metronome: &MockBillableMetricClient{
					CreateBillableMetricFn: func(ctx context.Context, reqData metronomeClient.CreateBillableMetricRequest) (*metronomeClient.CreateBillableMetricResponse, error) {
						expected := metronomeClient.CreateBillableMetricRequest{
							Name:            "name",
							AggregationType: metronomeClient.AggregationSum,
							AggregationKey:  "agg-key",
						}
						if diff := cmp.Diff(expected, reqData); diff != "" {
							t.Errorf("CreateBillableMetric mismatched: -want req, +got req: %s", diff)
						}

						return &metronomeClient.CreateBillableMetricResponse{
							Data: metronomeClient.BillableMetric{
								ID: "id",
							},
						}, nil
					},
				},
				mg: billableMetric(func(mg *v1alpha1.BillableMetric) {
					mg.Spec.ForProvider = v1alpha1.BillableMetricParameters{
						Name:            "name",
						AggregationType: v1alpha1.AggregationTypeSum,
						AggregationKey:  "agg-key",
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
