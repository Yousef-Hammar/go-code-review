package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/repository/memory"
)

var (
	ErrInvalidCode           = errors.New("invalid code")
	ErrNotFound              = errors.New("coupon not found")
	ErrInvalidDiscount       = errors.New("invalid discount")
	ErrInvalidMinBasketValue = errors.New("invalid min basket")
	ErrInvalidBasketValue    = errors.New("invalid basket value")
	ErrMinBasketValue        = errors.New("not sufficient basket value")
)

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) CreateCoupon(ctx context.Context, discount int, code string, minBasketValue int) error {
	if code == "" {
		return ErrInvalidCode
	}

	if discount < 0 || discount > 100 {
		return ErrInvalidDiscount
	}

	if minBasketValue < 0 {
		return ErrInvalidMinBasketValue
	}

	if _, err := s.repo.FindByCode(ctx, code); err == nil || !errors.Is(err, memory.ErrNotFound) {
		return ErrInvalidCode
	}

	coupon := domain.Coupon{
		ID:             uuid.NewString(),
		Code:           code,
		Discount:       discount,
		MinBasketValue: minBasketValue,
	}

	if err := s.repo.Save(ctx, coupon); err != nil {
		return err
	}
	return nil
}

func (s Service) GetCoupons(ctx context.Context, codes []string) ([]domain.Coupon, error) {
	coupons := make([]domain.Coupon, 0, len(codes))

	for _, code := range codes {
		coupon, err := s.repo.FindByCode(ctx, code)
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

func (s Service) ApplyCoupon(ctx context.Context, basket domain.Basket, code string) (*domain.Basket, error) {
	var (
		err error
	)

	if code == "" {
		return nil, ErrInvalidCode
	}

	if basket.Value <= 0 {
		return nil, ErrInvalidBasketValue
	}

	coupon, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		switch err {
		case memory.ErrNotFound:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	if basket.Value < coupon.Discount {
		return nil, ErrInvalidBasketValue
	}

	if basket.Value < coupon.MinBasketValue {
		return nil, ErrMinBasketValue
	}

	return &domain.Basket{
		Value:           basket.Value - coupon.Discount,
		AppliedDiscount: coupon.Discount,
	}, nil
}
