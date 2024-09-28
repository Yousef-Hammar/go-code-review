package memory

import (
	"errors"

	"coupon_service/internal/domain"
)

var ErrNotFound = errors.New("coupon not found")

type Repository struct {
	entries map[string]domain.Coupon
}

func New() *Repository {
	return &Repository{
		entries: make(map[string]domain.Coupon),
	}
}

func (r *Repository) FindByCode(code string) (*domain.Coupon, error) {
	coupon, ok := r.entries[code]
	if !ok {
		return nil, ErrNotFound
	}
	return &coupon, nil
}

func (r *Repository) Save(coupon domain.Coupon) error {
	r.entries[coupon.Code] = coupon
	return nil
}
