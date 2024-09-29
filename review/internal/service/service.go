package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	. "github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/repository/memory"
)

var (
	ErrInvalidCode           = errors.New("invalid code")
	ErrInvalidDiscount       = errors.New("invalid discount")
	ErrInvalidMinBasketValue = errors.New("invalid min basket")
)

type Repository interface {
	FindByCode(string) (*Coupon, error)
	Save(Coupon) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) CreateCoupon(discount int, code string, minBasketValue int) error {
	if code == "" {
		return ErrInvalidCode
	}

	if discount < 0 || discount > 100 {
		return ErrInvalidDiscount
	}

	if minBasketValue < 0 {
		return ErrInvalidMinBasketValue
	}

	if _, err := s.repo.FindByCode(code); err == nil || !errors.Is(err, memory.ErrNotFound) {
		return ErrInvalidCode
	}

	coupon := Coupon{
		ID:             uuid.NewString(),
		Code:           code,
		Discount:       discount,
		MinBasketValue: minBasketValue,
	}

	if err := s.repo.Save(coupon); err != nil {
		return err
	}
	return nil
}

func (s Service) GetCoupons(codes []string) ([]Coupon, error) {
	coupons := make([]Coupon, 0, len(codes))

	for _, code := range codes {
		coupon, err := s.repo.FindByCode(code)
		if err != nil {
			if errors.Is(err, memory.ErrNotFound) {
				continue
			}
			return nil, err
		}
		coupons = append(coupons, *coupon)
	}

	return coupons, nil
}

func (s Service) ApplyCoupon(basket Basket, code string) (b *Basket, e error) {
	b = &basket
	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	if b.Value > 0 {
		b.AppliedDiscount = coupon.Discount
		b.ApplicationSuccessful = true
	}
	if b.Value == 0 {
		return
	}

	return nil, fmt.Errorf("Tried to apply discount to negative value")
}
