package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/api"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/config"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/repository/memory"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/service"
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Coupon Service API Suite")
}

var _ = Describe("Coupon Service API", func() {
	var (
		router *gin.Engine
		app    *api.Application
		srv    service.Service
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()

		cfg := config.New()
		logger := zap.NewNop().Sugar()

		repo := memory.New()
		srv = service.New(repo)
		app = api.New(cfg, logger, srv)

		router.POST("/coupons", app.Create)
		router.GET("/coupons", app.Get)
		router.POST("/coupons/basket", app.Apply)
	})

	Describe("Creating a coupon", func() {
		Context("with valid input", func() {
			It("should create a new coupon and return 201", func() {
				body := api.CreateCouponReq{
					Code:           "test",
					Discount:       10,
					MinBasketValue: 100,
				}
				jsonBody, _ := json.Marshal(body)
				req, _ := http.NewRequest("POST", "/coupons", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusCreated))
			})
		})

		Context("with invalid input", func() {
			It("should return 400 for invalid discount", func() {
				body := api.CreateCouponReq{
					Code:           "test",
					Discount:       -10,
					MinBasketValue: 100,
				}
				jsonBody, _ := json.Marshal(body)
				req, _ := http.NewRequest("POST", "/coupons", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for invalid min basket", func() {
				body := api.CreateCouponReq{
					Code:           "test",
					Discount:       10,
					MinBasketValue: -10,
				}
				jsonBody, _ := json.Marshal(body)
				req, _ := http.NewRequest("POST", "/coupons", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("Getting coupons", func() {
		BeforeEach(func() {
			err := srv.CreateCoupon(nil, 10, "test", 100)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("with valid code", func() {
			It("should return the coupon details", func() {
				req, _ := http.NewRequest("GET", "/coupons?codes=test", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				expectedBody := `{"data":[{"code":"test","discount":10,"minBasketValue":100}]}`
				Expect(w.Body.String()).To(MatchJSON(expectedBody))
			})
		})

		Context("with no code specified", func() {
			It("should return 400", func() {
				req, _ := http.NewRequest("GET", "/coupons", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with non existing coupon for specified code", func() {
			It("should return 404", func() {
				req, _ := http.NewRequest("GET", "/coupons?codes=test2", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
