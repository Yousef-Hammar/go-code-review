package entity

import (
	"coupon_service/internal/domain"
)

type ApplicationRequest struct {
	Code   string
	Basket domain.Basket
}
