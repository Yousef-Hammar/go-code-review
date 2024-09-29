package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/repository/memory"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/service"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/service/internal/mocks"
)

func TestCreateCoupon(t *testing.T) {
	type args struct {
		code           string
		discount       int
		minBasketValue int
	}

	type testCase struct {
		name        string
		args        args
		setupMocks  func(*mocks.Repository, args)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Successful coupon creation",
			args: args{code: "test", discount: 10, minBasketValue: 5},
			setupMocks: func(repo *mocks.Repository, args args) {
				repo.On("FindByCode", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), args.code).
					Return(nil, memory.ErrNotFound).
					Once()
				repo.On("Save", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), mock.MatchedBy(func(coupon domain.Coupon) bool {
					return coupon.ID != "" &&
						coupon.Code == args.code &&
						coupon.Discount == args.discount &&
						coupon.MinBasketValue == args.minBasketValue
				})).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:        "Empty coupon code",
			args:        args{code: "", discount: 10, minBasketValue: 5},
			setupMocks:  func(repo *mocks.Repository, args args) {},
			expectedErr: service.ErrInvalidCode,
		},
		{
			name: "Duplicated coupon code",
			args: args{code: "test", discount: 10, minBasketValue: 5},
			setupMocks: func(repo *mocks.Repository, args args) {
				repo.On("FindByCode", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), args.code).
					Return(&domain.Coupon{}, nil).
					Once()
			},
			expectedErr: service.ErrInvalidCode,
		},
		{
			name:        "Negative discount value",
			args:        args{code: "test", discount: -1, minBasketValue: 5},
			setupMocks:  func(repo *mocks.Repository, args args) {},
			expectedErr: service.ErrInvalidDiscount,
		},
		{
			name:        "Discount value greater than 100",
			args:        args{code: "test", discount: 200, minBasketValue: 5},
			setupMocks:  func(repo *mocks.Repository, args args) {},
			expectedErr: service.ErrInvalidDiscount,
		},
		{
			name:        "Negative minimum basket value",
			args:        args{code: "test", discount: 10, minBasketValue: -1},
			setupMocks:  func(repo *mocks.Repository, args args) {},
			expectedErr: service.ErrInvalidMinBasketValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			tc.setupMocks(repo, tc.args)
			defer repo.AssertExpectations(t)

			srv := service.New(repo)
			ctx := context.Background()

			err := srv.CreateCoupon(ctx, tc.args.discount, tc.args.code, tc.args.minBasketValue)
			if tc.expectedErr != nil {
				assert.Error(t, err, "expected error to be %v, got: %v", tc.expectedErr, err)
				assert.IsType(t, tc.expectedErr, err, "expected error %T, got: %T", tc.expectedErr, err)
				return
			}

			assert.NoError(t, err, "expected error nil, got: %v", err)
		})

	}
}

func TestGetCoupon(t *testing.T) {
	type testCase struct {
		name        string
		codes       []string
		setupMocks  func(*mocks.Repository)
		want        []domain.Coupon
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "Successful coupons retrieval",
			codes: []string{"test1", "test2"},
			setupMocks: func(repo *mocks.Repository) {
				repo.On("FindByCode", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), "test1").
					Return(&domain.Coupon{ID: "id1", Code: "test1", Discount: 10, MinBasketValue: 0}, nil).
					Once()
				repo.On("FindByCode", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), "test2").
					Return(&domain.Coupon{ID: "id2", Code: "test2", Discount: 10, MinBasketValue: 0}, nil).
					Once()
			},
			want: []domain.Coupon{
				{ID: "id1", Code: "test1", Discount: 10, MinBasketValue: 0},
				{ID: "id2", Code: "test2", Discount: 10, MinBasketValue: 0},
			},
			expectedErr: nil,
		},
		{
			name:  "Successful coupons retrieval with no existing coupon",
			codes: []string{"test1", "test2"},
			setupMocks: func(repo *mocks.Repository) {
				repo.On("FindByCode", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), "test1").
					Return(&domain.Coupon{ID: "id1", Code: "test1", Discount: 10, MinBasketValue: 0}, nil).
					Once()
				repo.On("FindByCode", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), "test2").
					Return(nil, memory.ErrNotFound).Once().
					Once()
			},
			want: []domain.Coupon{
				{ID: "id1", Code: "test1", Discount: 10, MinBasketValue: 0},
			},
			expectedErr: nil,
		},
		{
			name:  "Error during coupon retrieval",
			codes: []string{"test1", "test2"},
			setupMocks: func(repo *mocks.Repository) {
				repo.On("FindByCode", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), "test1").
					Return(nil, errors.New("fatal error")).
					Once()
			},
			want:        nil,
			expectedErr: errors.New("fatal error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			tc.setupMocks(repo)
			defer repo.AssertExpectations(t)

			srv := service.New(repo)
			ctx := context.Background()

			got, err := srv.GetCoupons(ctx, tc.codes)
			if tc.expectedErr != nil {
				assert.Error(t, err, "expected error to be %v, got: %v", tc.expectedErr, err)
				assert.IsType(t, tc.expectedErr, err, "expected error %v, got: %v", tc.expectedErr, err)
				return
			}
			assert.EqualValues(t, tc.want, got, "expected coupons slice to be of length %d, "+
				"got slice of length %d", len(tc.want), len(got))
		})
	}
}

func TestApplyCoupon(t *testing.T) {
	type args struct {
		code   string
		basket domain.Basket
	}
	type testCase struct {
		name        string
		args        args
		setupMocks  func(*mocks.Repository, string)
		want        *domain.Basket
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Successful coupon application",
			args: args{
				code:   "test1",
				basket: domain.Basket{Value: 50},
			},
			setupMocks: func(repo *mocks.Repository, code string) {
				repo.On("FindByCode", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), code).Return(&domain.Coupon{
					ID:             "id1",
					Code:           code,
					Discount:       10,
					MinBasketValue: 20,
				}, nil).Once()
			},
			want: &domain.Basket{
				Value:                 45,
				AppliedDiscount:       10,
				ApplicationSuccessful: true,
			},
			expectedErr: nil,
		},
		{
			name: "Basket with negative value",
			args: args{
				code:   "test1",
				basket: domain.Basket{Value: -50},
			},
			setupMocks:  func(repo *mocks.Repository, code string) {},
			want:        nil,
			expectedErr: service.ErrInvalidBasketValue,
		},
		{
			name: "Basket with value less than minimum required",
			args: args{
				code:   "test1",
				basket: domain.Basket{Value: 10},
			},
			setupMocks: func(repo *mocks.Repository, code string) {
				repo.On("FindByCode", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), code).Return(&domain.Coupon{
					ID:             "id1",
					Code:           code,
					Discount:       10,
					MinBasketValue: 20,
				}, nil).Once()
			},
			want:        nil,
			expectedErr: service.ErrMinBasketValue,
		},
		{
			name: "Invalid coupon code",
			args: args{
				code:   "test1",
				basket: domain.Basket{Value: 10},
			},
			setupMocks: func(repo *mocks.Repository, code string) {
				repo.On("FindByCode", mock.MatchedBy(func(ctx context.Context) bool {
					return true
				}), code).Return(nil, service.ErrInvalidCode).Once()
			},
			want:        nil,
			expectedErr: service.ErrInvalidCode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			tc.setupMocks(repo, tc.args.code)
			defer repo.AssertExpectations(t)

			srv := service.New(repo)
			ctx := context.Background()

			got, err := srv.ApplyCoupon(ctx, tc.args.basket, tc.args.code)
			if tc.expectedErr != nil {
				assert.Error(t, err, "expected error to be %v, got: %v", tc.expectedErr, err)
				assert.IsType(t, tc.expectedErr, err, "expected error %v, got: %v", tc.expectedErr, err)
				return
			}
			assert.EqualValues(t, tc.want, got, "expected basket to be %+v, got: %+v", tc.want, got)
		})
	}
}
