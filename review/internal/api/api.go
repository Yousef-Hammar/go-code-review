package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Application struct {
	config  Config
	logger  *zap.SugaredLogger
	service Service
}

func NewApplication(config Config, logger *zap.SugaredLogger, service Service) *Application {
	return &Application{
		config:  config,
		logger:  logger,
		service: service,
	}
}

func (app *Application) Mount() http.Handler {
	router := gin.New()

	v1 := router.Group("/v1")

	v1.POST("/coupons", func(context *gin.Context) {
		app.logger.Infof("hello world from CreateCoupon")
	})
	v1.GET("/coupons", func(context *gin.Context) {
		app.logger.Infof("hello world from GetCoupons")
	})
	v1.POST("/coupons/basket", func(context *gin.Context) {
		app.logger.Infof("hello world from ApplyCoupon")
	})

	return router
}

func (app *Application) Run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("start http server on", "addr", srv.Addr)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		app.logger.Fatalw("start http server failed", "error", err)
		return err
	}

	return nil
}
