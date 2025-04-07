package product

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

	"github.com/redbackthomson/provider-metronome/apis/product/v1alpha1"
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

type productModifier func(mg *v1alpha1.Product)

func product(rm ...productModifier) *v1alpha1.Product {
	r := &v1alpha1.Product{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testResourceName,
			Namespace: testNamespace,
		},
		Spec: v1alpha1.ProductSpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{
					Name: providerConfigName,
				},
			},
			ForProvider: v1alpha1.ProductParameters{},
		},
		Status: v1alpha1.ProductStatus{},
	}

	meta.SetExternalName(r, "external-name")

	for _, m := range rm {
		m(r)
	}

	return r
}

type notProductResource struct {
	resource.Managed
}

type MockProductClient struct {
	ArchiveProductFn func(ctx context.Context, reqData metronomeClient.ArchiveProductRequest) (*metronomeClient.ArchiveProductResponse, error)
	CreateProductFn  func(ctx context.Context, reqData metronomeClient.CreateProductRequest) (*metronomeClient.CreateProductResponse, error)
	GetProductFn     func(ctx context.Context, reqData metronomeClient.GetProductRequest) (*metronomeClient.GetProductResponse, error)
	ListProductFn    func(ctx context.Context, reqData metronomeClient.ListProductsRequest, nextPage string) (*metronomeClient.ListProductsResponse, error)
	UpdateProductFn  func(ctx context.Context, reqData metronomeClient.UpdateProductRequest) (*metronomeClient.UpdateProductResponse, error)
}

func (m *MockProductClient) ArchiveProduct(ctx context.Context, reqData metronomeClient.ArchiveProductRequest) (*metronomeClient.ArchiveProductResponse, error) {
	return m.ArchiveProductFn(ctx, reqData)
}

func (m *MockProductClient) CreateProduct(ctx context.Context, reqData metronomeClient.CreateProductRequest) (*metronomeClient.CreateProductResponse, error) {
	return m.CreateProductFn(ctx, reqData)
}

func (m *MockProductClient) GetProduct(ctx context.Context, reqData metronomeClient.GetProductRequest) (*metronomeClient.GetProductResponse, error) {
	return m.GetProductFn(ctx, reqData)
}

func (m *MockProductClient) ListProduct(ctx context.Context, reqData metronomeClient.ListProductsRequest, nextPage string) (*metronomeClient.ListProductsResponse, error) {
	return m.ListProductFn(ctx, reqData, nextPage)
}

func (m *MockProductClient) UpdateProduct(ctx context.Context, reqData metronomeClient.UpdateProductRequest) (*metronomeClient.UpdateProductResponse, error) {
	return m.UpdateProductFn(ctx, reqData)
}

var _ (metronomeClient.ProductClient) = (*MockProductClient)(nil)

func Test_External_Observe(t *testing.T) {
	type args struct {
		metronome metronomeClient.ProductClient
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
		"NotProductResource": {
			args: args{
				mg: notProductResource{},
			},
			want: want{
				err: errors.New(errNotProduct),
			},
		},
		"FailedToGetProduct": {
			args: args{
				metronome: &MockProductClient{
					GetProductFn: func(ctx context.Context, reqData metronomeClient.GetProductRequest) (*metronomeClient.GetProductResponse, error) {
						return nil, errBoom
					},
				},
				mg: product(),
			},
			want: want{
				err: errors.Wrap(errBoom, errGetProduct),
			},
		},
		"ArchivedIsDeleted": {
			args: args{
				metronome: &MockProductClient{
					GetProductFn: func(ctx context.Context, reqData metronomeClient.GetProductRequest) (*metronomeClient.GetProductResponse, error) {
						return &metronomeClient.GetProductResponse{
							Data: metronomeClient.Product{
								ArchivedAt: "archived-at",
							},
						}, nil
					},
				},
				mg: product(),
			},
			want: want{
				out: managed.ExternalObservation{ResourceExists: false},
			},
		},
		"NotUpToDate": {
			args: args{
				metronome: &MockProductClient{
					GetProductFn: func(ctx context.Context, reqData metronomeClient.GetProductRequest) (*metronomeClient.GetProductResponse, error) {
						return &metronomeClient.GetProductResponse{
							Data: metronomeClient.Product{
								ID: "id1", Current: metronomeClient.ProductDetails{Name: "name"},
							},
						}, nil
					},
				},
				mg: product(func(mg *v1alpha1.Product) {
					mg.Spec.ForProvider = v1alpha1.ProductParameters{
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
				metronome: &MockProductClient{
					GetProductFn: func(ctx context.Context, reqData metronomeClient.GetProductRequest) (*metronomeClient.GetProductResponse, error) {
						return &metronomeClient.GetProductResponse{
							Data: metronomeClient.Product{
								ID: "id1", Current: metronomeClient.ProductDetails{Name: "name"},
							},
						}, nil
					},
				},
				mg: product(func(mg *v1alpha1.Product) {
					meta.SetExternalName(mg, "id1")
					mg.Spec.ForProvider = v1alpha1.ProductParameters{
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
		metronome metronomeClient.ProductClient
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
		"NotProductResource": {
			args: args{
				mg: notProductResource{},
			},
			want: want{
				err: errors.New(errNotProduct),
			},
		},
		"FailedToCreateProduct": {
			args: args{
				metronome: &MockProductClient{
					CreateProductFn: func(ctx context.Context, reqData metronomeClient.CreateProductRequest) (*metronomeClient.CreateProductResponse, error) {
						return nil, errBoom
					},
				},
				mg: product(),
			},
			want: want{
				err: errors.Wrap(errBoom, errCreateProduct),
			},
		},
		"Success": {
			args: args{
				metronome: &MockProductClient{
					CreateProductFn: func(ctx context.Context, reqData metronomeClient.CreateProductRequest) (*metronomeClient.CreateProductResponse, error) {
						expected := metronomeClient.CreateProductRequest{
							Name: "name",
						}
						if diff := cmp.Diff(expected, reqData); diff != "" {
							t.Errorf("CreateProduct mismatched: -want req, +got req: %s", diff)
						}

						return &metronomeClient.CreateProductResponse{
							Data: metronomeClient.IDOnly{ID: "id"},
						}, nil
					},
				},
				mg: product(func(mg *v1alpha1.Product) {
					mg.Spec.ForProvider = v1alpha1.ProductParameters{
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

func Test_External_Delete(t *testing.T) {
	type args struct {
		metronome metronomeClient.ProductClient
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
		"NotProductResource": {
			args: args{
				mg: notProductResource{},
			},
			want: want{
				err: errors.New(errNotProduct),
			},
		},
		"FailedToDeleteProduct": {
			args: args{
				metronome: &MockProductClient{
					ArchiveProductFn: func(ctx context.Context, reqData metronomeClient.ArchiveProductRequest) (*metronomeClient.ArchiveProductResponse, error) {
						return nil, errBoom
					},
				},
				mg: product(),
			},
			want: want{
				err: errors.Wrap(errBoom, errArchiveProduct),
			},
		},
		"Success": {
			args: args{
				metronome: &MockProductClient{
					ArchiveProductFn: func(ctx context.Context, reqData metronomeClient.ArchiveProductRequest) (*metronomeClient.ArchiveProductResponse, error) {
						expected := metronomeClient.ArchiveProductRequest{
							ProductID: "product-id",
						}
						if diff := cmp.Diff(expected, reqData); diff != "" {
							t.Errorf("ArchiveProduct mismatched: -want req, +got req: %s", diff)
						}

						return nil, nil
					},
				},
				mg: product(func(mg *v1alpha1.Product) {
					meta.SetExternalName(mg, "product-id")
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
