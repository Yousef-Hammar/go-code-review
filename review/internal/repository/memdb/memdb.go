package memdb

import (
	"fmt"

	"coupon_service/internal/domain"
)

type Config struct{}

type repository interface {
	FindByCode(string) (*domain.Coupon, error)
	Save(domain.Coupon) error
}

type Repository struct {
	entries map[string]domain.Coupon
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) FindByCode(code string) (*domain.Coupon, error) {
	coupon, ok := r.entries[code]
	if !ok {
		return nil, fmt.Errorf("Coupon not found")
	}
	return &coupon, nil
}

func (r *Repository) Save(coupon domain.Coupon) error {
	r.entries[coupon.Code] = coupon
	return nil
}
