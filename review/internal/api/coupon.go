package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/service"
)

type CreateCouponReq struct {
	Code           string `json:"code" binding:"required"`
	Discount       int    `json:"discount" binding:"required"`
	MinBasketValue int    `json:"minBasketValue" binding:"required"`
}

func (app *Application) CreateCoupon(c *gin.Context) {
	var (
		body CreateCouponReq
	)

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := app.service.CreateCoupon(c.Request.Context(), body.Discount, body.Code, body.MinBasketValue)
	if err != nil {
		switch err {
		case service.ErrInvalidCode, service.ErrInvalidDiscount, service.ErrInvalidMinBasketValue:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.Status(http.StatusCreated)
}

type Coupon struct {
	Code           string `json:"code"`
	Discount       int    `json:"discount"`
	MinBasketValue int    `json:"minBasketValue"`
}

func (app *Application) Get(c *gin.Context) {
	var (
		resp []Coupon
	)

	rawCodes := c.Query("codes")

	if rawCodes == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no code specified"})
		return
	}

	codes := strings.Split(rawCodes, ",")

	coupons, err := app.service.GetCoupons(c.Request.Context(), codes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp = make([]Coupon, 0, len(coupons))
	for _, coupon := range coupons {
		resp = append(resp, Coupon{
			Code:           coupon.Code,
			Discount:       coupon.Discount,
			MinBasketValue: coupon.MinBasketValue,
		})
	}

	c.JSON(http.StatusOK, resp)
}

type Basket struct {
	Value                 int  `json:"value"`
	AppliedDiscount       int  `json:"appliedDiscount"`
	ApplicationSuccessful bool `json:"applicationSuccessful"`
}

type ApplyReq struct {
	Basket Basket `json:"basket"`
	Code   string `json:"code"`
}

func (app *Application) Apply(c *gin.Context) {
	var body ApplyReq

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	basket := &domain.Basket{
		Value: body.Basket.Value,
	}

	basket, err := app.service.ApplyCoupon(c.Request.Context(), *basket, body.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Basket{
		Value:                 basket.Value,
		AppliedDiscount:       basket.AppliedDiscount,
		ApplicationSuccessful: basket.ApplicationSuccessful,
	})
}
