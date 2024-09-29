package memory

import (
	"errors"
	"sync"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
)

var ErrNotFound = errors.New("coupon not found")

type Repository struct {
	entries map[string]domain.Coupon
	mu      *sync.Mutex
}

func New() *Repository {
	return &Repository{
		entries: make(map[string]domain.Coupon),
		mu:      &sync.Mutex{},
	}
}

func (r *Repository) FindByCode(code string) (*domain.Coupon, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if coupon, ok := r.entries[code]; ok {
		return &coupon, nil
	}
	return nil, ErrNotFound
}

func (r *Repository) Save(coupon domain.Coupon) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.entries[coupon.Code] = coupon
	return nil
}
