package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/repository/memory"
)

var (
	ErrInvalidCode           = errors.New("invalid code")
	ErrInvalidDiscount       = errors.New("invalid discount")
	ErrInvalidMinBasketValue = errors.New("invalid min basket")
	ErrInvalidBasketValue    = errors.New("invalid basket value")
	ErrMinBasketValue        = errors.New("not sufficient basket value")
)

type Repository interface {
	FindByCode(string) (*domain.Coupon, error)
	Save(domain.Coupon) error
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

	coupon := domain.Coupon{
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

func (s Service) GetCoupons(codes []string) ([]domain.Coupon, error) {
	coupons := make([]domain.Coupon, 0, len(codes))

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

func (s Service) ApplyCoupon(basket domain.Basket, code string) (*domain.Basket, error) {
	var (
		err error
	)

	if basket.Value <= 0 {
		return nil, ErrInvalidBasketValue
	}

	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	if basket.Value < coupon.MinBasketValue {
		return nil, ErrMinBasketValue
	}

	discountAmount := (basket.Value * coupon.Discount) / 100

	return &domain.Basket{
		Value:                 basket.Value - discountAmount,
		AppliedDiscount:       coupon.Discount,
		ApplicationSuccessful: true,
	}, nil
}
