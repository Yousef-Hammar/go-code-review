package memory_test

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/repository/memory"
)

func TestFindByCode(t *testing.T) {
	if os.Getenv("LONG") != "" {
		t.Skip("Skipping TestFindByCode in long mode.")
	}

	type testCase struct {
		name        string
		code        string
		expectedErr error
		want        *domain.Coupon
	}

	testCases := []testCase{
		{
			name:        "Coupon found",
			code:        "test",
			expectedErr: nil,
			want: &domain.Coupon{
				ID:             "test",
				Code:           "test",
				Discount:       0,
				MinBasketValue: 0,
			},
		},
		{
			name:        "Coupon not found",
			code:        "not found",
			expectedErr: memory.ErrNotFound,
			want:        nil,
		},
		{
			name:        "Empty coupon code",
			code:        "",
			expectedErr: memory.ErrNotFound,
			want:        nil,
		},
	}

	ctx := context.Background()
	repo := memory.New()
	_ = repo.Save(ctx, domain.Coupon{
		ID:             "test",
		Code:           "test",
		Discount:       0,
		MinBasketValue: 0,
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			coupon, err := repo.FindByCode(ctx, tc.code)
			if tc.expectedErr != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("expected err to be %v, got %v", tc.expectedErr, err)
					return
				}
			}
			if !reflect.DeepEqual(tc.want, coupon) {
				t.Errorf("expected coupon to be %v, got %v", tc.want, coupon)
			}
		})
	}
}

func TestSave(t *testing.T) {
	if os.Getenv("LONG") != "" {
		t.Skip("Skipping TestSave in long mode.")
	}

	type testCase struct {
		name        string
		coupon      domain.Coupon
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Successful save",
			coupon: domain.Coupon{
				ID:             "test",
				Code:           "test",
				Discount:       0,
				MinBasketValue: 0,
			},
			expectedErr: nil,
		},
	}

	repo := memory.New()
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Save(ctx, tc.coupon)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected err to be %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
