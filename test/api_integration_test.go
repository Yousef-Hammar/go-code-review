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
		router http.Handler
		app    *api.Application
		srv    service.Service
	)

	BeforeEach(func() {
		cfg := config.New()
		logger := zap.NewNop().Sugar()

		repo := memory.New()
		srv = service.New(repo)
		app = api.New(cfg, logger, srv)

		router = app.Mount(gin.TestMode)
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
				req, _ := http.NewRequest(http.MethodPost, "/v1/coupons", bytes.NewBuffer(jsonBody))
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
				req, _ := http.NewRequest(http.MethodPost, "/v1/coupons", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(ContainSubstring("invalid discount"))
			})

			It("should return 400 for invalid min basket", func() {
				body := api.CreateCouponReq{
					Code:           "test",
					Discount:       10,
					MinBasketValue: -10,
				}
				jsonBody, _ := json.Marshal(body)
				req, _ := http.NewRequest(http.MethodPost, "/v1/coupons", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(ContainSubstring("invalid min basket"))
			})

			It("should return 400 for empty code", func() {
				body := api.CreateCouponReq{
					Code:           "",
					Discount:       10,
					MinBasketValue: 100,
				}
				jsonBody, _ := json.Marshal(body)
				req, _ := http.NewRequest(http.MethodPost, "/v1/coupons", bytes.NewBuffer(jsonBody))
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
				req, _ := http.NewRequest(http.MethodGet, "/v1/coupons?codes=test", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				expectedBody := `{"data":[{"code":"test","discount":10,"minBasketValue":100}]}`
				Expect(w.Body.String()).To(MatchJSON(expectedBody))
			})
		})

		Context("with no code specified", func() {
			It("should return 400", func() {
				req, _ := http.NewRequest(http.MethodGet, "/v1/coupons", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with non-existing coupon for specified code", func() {
			It("should return 404", func() {
				req, _ := http.NewRequest(http.MethodGet, "/v1/coupons?codes=test2", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with multiple codes", func() {
			BeforeEach(func() {
				err := srv.CreateCoupon(nil, 20, "test2", 200)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return all existing coupons", func() {
				req, _ := http.NewRequest(http.MethodGet, "/v1/coupons?codes=test,test2,nonexistent", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				expectedBody := `{"data":[{"code":"test","discount":10,"minBasketValue":100},{"code":"test2","discount":20,"minBasketValue":200}]}`
				Expect(w.Body.String()).To(MatchJSON(expectedBody))
			})
		})
	})

	Describe("Applying a coupon", func() {
		BeforeEach(func() {
			err := srv.CreateCoupon(nil, 10, "test", 100)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("with valid basket and coupon", func() {
			It("should apply the discount and return updated basket", func() {
				body := api.ApplyReq{
					Basket: api.Basket{Value: 200},
					Code:   "test",
				}
				jsonBody, _ := json.Marshal(body)
				req, _ := http.NewRequest(http.MethodPost, "/v1/coupons/basket", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				expectedBody := `{"data":{"value":190,"appliedDiscount":10}}`
				Expect(w.Body.String()).To(MatchJSON(expectedBody))
			})
		})

		Context("with basket value below minimum", func() {
			It("should return 400", func() {
				body := api.ApplyReq{
					Basket: api.Basket{Value: 50},
					Code:   "test",
				}
				jsonBody, _ := json.Marshal(body)
				req, _ := http.NewRequest(http.MethodPost, "/v1/coupons/basket", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with non-existent coupon code", func() {
			It("should return 404", func() {
				body := api.ApplyReq{
					Basket: api.Basket{Value: 200},
					Code:   "nonexistent",
				}
				jsonBody, _ := json.Marshal(body)
				req, _ := http.NewRequest(http.MethodPost, "/v1/coupons/basket", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
