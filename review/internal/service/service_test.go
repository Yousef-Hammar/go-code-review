package service_test

import (
	"reflect"
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
				repo.On("FindByCode", args.code).Return(nil, memory.ErrNotFound).Once()
				repo.On("Save", mock.MatchedBy(func(coupon domain.Coupon) bool {
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
				repo.On("FindByCode", args.code).Return(&domain.Coupon{}, nil).Once()
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

			err := srv.CreateCoupon(tc.args.discount, tc.args.code, tc.args.minBasketValue)
			if tc.expectedErr != nil {
				assert.Error(t, err, "expected error to be %v, got: %v", tc.expectedErr, err)
				assert.IsType(t, tc.expectedErr, err, "expected error %T, got: %T", tc.expectedErr, err)
				return
			}

			assert.NoError(t, err, "expected error nil, got: %v", err)
		})

	}
}

func TestService_ApplyCoupon(t *testing.T) {
	type fields struct {
		repo service.Repository
	}
	type args struct {
		basket domain.Basket
		code   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantB   *domain.Basket
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.New(tt.fields.repo)
			gotB, err := s.ApplyCoupon(tt.args.basket, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyCoupon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("ApplyCoupon() gotB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
