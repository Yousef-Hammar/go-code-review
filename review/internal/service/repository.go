package service

import "github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"

type Repository interface {
	FindByCode(string) (*domain.Coupon, error)
	Save(domain.Coupon) error
}
