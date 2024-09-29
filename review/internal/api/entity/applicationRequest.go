package entity

import (
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
)

type ApplicationRequest struct {
	Code   string
	Basket domain.Basket
}
