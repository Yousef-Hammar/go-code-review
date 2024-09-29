package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
