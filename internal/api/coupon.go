package api

import (
	"errors"
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

func (app *Application) Create(c *gin.Context) {
	var (
		body CreateCouponReq
	)

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		app.logger.Errorw("error occurred while binding body", "error", err)
		app.writeJSONError(c, http.StatusBadRequest, err)
		return
	}

	err := app.service.CreateCoupon(c.Request.Context(), body.Discount, body.Code, body.MinBasketValue)
	if err != nil {
		app.logger.Errorw("error occurred while creating coupon", "error", err)
		switch err {
		case service.ErrInvalidCode, service.ErrInvalidDiscount, service.ErrInvalidMinBasketValue:
			app.writeJSONError(c, http.StatusBadRequest, err)
			return
		default:
			app.writeJSONError(c, http.StatusInternalServerError, err)
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
		app.logger.Errorw("error occurred while getting coupons, missing codes")
		app.writeJSONError(c, http.StatusBadRequest, errors.New("no code specified"))
		return
	}

	codes := strings.Split(rawCodes, ",")

	coupons, err := app.service.GetCoupons(c.Request.Context(), codes)
	if err != nil {
		app.logger.Errorw("error occurred while getting coupons", "error", err)
		app.writeJSONError(c, http.StatusInternalServerError, err)
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

	if len(resp) == 0 {
		app.logger.Debug("no coupons found", "codes", rawCodes)
		app.writeJSONError(c, http.StatusNotFound, errors.New("no coupons found"))
		return
	}

	app.writeJSONResponse(c, http.StatusOK, resp)
}

type Basket struct {
	Value           int `json:"value" binding:"required"`
	AppliedDiscount int `json:"appliedDiscount"`
}

type ApplyReq struct {
	Basket Basket `json:"basket" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

func (app *Application) Apply(c *gin.Context) {
	var body ApplyReq

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		app.logger.Errorw("error occurred while binding body", "error", err)
		app.writeJSONError(c, http.StatusBadRequest, err)
		return
	}

	basket := &domain.Basket{
		Value: body.Basket.Value,
	}

	basket, err := app.service.ApplyCoupon(c.Request.Context(), *basket, body.Code)
	if err != nil {
		app.logger.Errorw("error occurred while applying coupon", "error", err)
		switch err {
		case service.ErrInvalidCode, service.ErrInvalidBasketValue, service.ErrMinBasketValue, service.ErrNotFound:
			app.writeJSONError(c, http.StatusBadRequest, err)
			return
		default:
			app.writeJSONError(c, http.StatusInternalServerError, err)
			return
		}
	}

	app.writeJSONResponse(c, http.StatusOK, Basket{
		Value:           basket.Value,
		AppliedDiscount: basket.AppliedDiscount,
	})
}
