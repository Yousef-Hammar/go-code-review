package api

import (
	"context"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
)

type Service interface {
	CreateCoupon(context.Context, int, string, int) error
	GetCoupons(context.Context, []string) ([]domain.Coupon, error)
	ApplyCoupon(context.Context, domain.Basket, string) (*domain.Basket, error)
}
