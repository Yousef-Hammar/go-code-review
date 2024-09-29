package service

import (
	"context"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
)

type Repository interface {
	FindByCode(context.Context, string) (*domain.Coupon, error)
	Save(context.Context, domain.Coupon) error
}
