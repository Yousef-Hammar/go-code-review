package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type config struct {
	Addr string
}

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

	v1.POST("/coupons", app.Create)
	v1.GET("/coupons", app.Get)
	v1.POST("/coupons/basket", app.Apply)

	return router
}

func (app *Application) Run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s", app.config.Addr),
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("start http server on", "addr", srv.Addr)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		app.logger.Fatalw("start http server failed", "error", err)
		return err
	}

	return nil
}
