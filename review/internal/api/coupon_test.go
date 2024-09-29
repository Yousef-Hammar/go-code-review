package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/api"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/api/internal/mocks"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/service"
)

func newTestApplication(t *testing.T, srv *mocks.Service) *api.Application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	return api.NewApplication(api.Config{}, logger, srv)
}

func TestCreate(t *testing.T) {
	type testCase struct {
		name      string
		body      *api.CreateCouponReq
		setupMock func(srv *mocks.Service, args *api.CreateCouponReq)
		want      int
	}

	tests := []testCase{
		{
			name: "Successful coupon creation",
			body: &api.CreateCouponReq{
				Code:           "test",
				Discount:       10,
				MinBasketValue: 20,
			},
			setupMock: func(srv *mocks.Service, args *api.CreateCouponReq) {
				srv.On("CreateCoupon",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					args.Discount, args.Code, args.MinBasketValue).
					Return(nil).
					Once()

			},
			want: http.StatusCreated,
		},
		{
			name:      "Invalid body",
			body:      nil,
			setupMock: func(srv *mocks.Service, args *api.CreateCouponReq) {},
			want:      http.StatusBadRequest,
		},
		{
			name: "Negative discount",
			body: &api.CreateCouponReq{
				Code:           "test",
				Discount:       -10,
				MinBasketValue: 20,
			},
			setupMock: func(srv *mocks.Service, args *api.CreateCouponReq) {
				srv.On("CreateCoupon", mock.MatchedBy(func(_ context.Context) bool { return true }),
					args.Discount, args.Code, args.MinBasketValue).
					Return(service.ErrInvalidDiscount).
					Once()
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Negative minimum basket value",
			body: &api.CreateCouponReq{
				Code:           "test",
				Discount:       10,
				MinBasketValue: -20,
			},
			setupMock: func(srv *mocks.Service, args *api.CreateCouponReq) {
				srv.On("CreateCoupon", mock.MatchedBy(func(_ context.Context) bool { return true }),
					args.Discount, args.Code, args.MinBasketValue).
					Return(service.ErrInvalidMinBasketValue).
					Once()
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Internal server error",
			body: &api.CreateCouponReq{
				Code:           "test",
				Discount:       20,
				MinBasketValue: 40,
			},
			setupMock: func(srv *mocks.Service, args *api.CreateCouponReq) {
				srv.On("CreateCoupon", mock.MatchedBy(func(_ context.Context) bool { return true }),
					args.Discount, args.Code, args.MinBasketValue).
					Return(errors.New("error")).
					Once()
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := mocks.NewService(t)
			tc.setupMock(srv, tc.body)
			defer srv.AssertExpectations(t)

			app := newTestApplication(t, srv)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/v1/coupons", app.CreateCoupon)
			httptest.NewServer(router)

			var buff bytes.Buffer
			err := json.NewEncoder(&buff).Encode(tc.body)
			require.NoErrorf(t, err, "error encoding request %v", err)

			req := httptest.NewRequest(http.MethodPost, "/v1/coupons", strings.NewReader(buff.String()))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.want, w.Code)
		})
	}
}

func TestGet(t *testing.T) {
	type coupon struct {
		Code           string `json:"code"`
		Discount       int    `json:"discount"`
		MinBasketValue int    `json:"minBasketValue"`
	}

	type testCase struct {
		name           string
		codes          []string
		setupMock      func(*mocks.Service, []string)
		wantStatusCode int
		want           []coupon
	}

	tests := []testCase{
		{
			name:  "Successful retrieval",
			codes: []string{"test", "test2"},
			setupMock: func(srv *mocks.Service, codes []string) {
				srv.On("GetCoupons", mock.MatchedBy(func(_ context.Context) bool { return true }), codes).
					Return([]domain.Coupon{
						{
							ID:             "id1",
							Code:           "test",
							Discount:       10,
							MinBasketValue: 20,
						},
						{
							ID:             "id2",
							Code:           "test2",
							Discount:       30,
							MinBasketValue: 500,
						},
					}, nil).
					Once()
			},
			wantStatusCode: http.StatusOK,
			want: []coupon{
				{Code: "test", Discount: 10, MinBasketValue: 20},
				{Code: "test2", Discount: 30, MinBasketValue: 500},
			},
		},
		{
			name:           "Empty codes",
			codes:          []string{},
			setupMock:      func(srv *mocks.Service, codes []string) {},
			wantStatusCode: http.StatusBadRequest,
			want:           []coupon{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := mocks.NewService(t)
			tc.setupMock(srv, tc.codes)
			defer srv.AssertExpectations(t)

			app := newTestApplication(t, srv)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.GET("/v1/coupons", app.Get)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/coupons?codes=%s", strings.Join(tc.codes, ",")), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatusCode, w.Code, "expected status code %d, got: %d", tc.wantStatusCode, w.Code)
			if tc.wantStatusCode == http.StatusOK {
				var resp []coupon
				require.NoError(t, json.NewDecoder(w.Body).Decode(&resp), "error decoding response body")

				assert.Equal(t, tc.want, resp, "expected %+v, got: %+v", tc.want, resp)
			}
		})
	}

}
