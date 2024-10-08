package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/config"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

type Application struct {
	config  config.Config
	logger  *zap.SugaredLogger
	service Service
}

func New(config config.Config, logger *zap.SugaredLogger, service Service) *Application {
	return &Application{
		config:  config,
		logger:  logger,
		service: service,
	}
}

func (app *Application) requestLoggerMiddleware(c *gin.Context) {
	var requestBody []byte
	var responseBody bytes.Buffer
	now := time.Now().UTC()

	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			app.logger.Error("Failed to read request body", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		requestBody = bodyBytes
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	c.Writer = &responseWriter{
		ResponseWriter: c.Writer,
		body:           &responseBody,
	}

	c.Next()

	latency := time.Since(now)

	app.logger.Info("Request",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("body", string(requestBody)),
		zap.String("clientIP", c.ClientIP()),
		zap.Int("status", c.Writer.Status()),
		zap.String("response", responseBody.String()),
		zap.Duration("latency", latency),
	)

	c.Request.Body = io.NopCloser(bytes.NewReader(requestBody))
}

func (app *Application) Mount(mode string) http.Handler {
	gin.SetMode(mode)
	router := gin.New()
	router.Use(app.requestLoggerMiddleware)
	router.Use(cors.Default())

	v1 := router.Group("/v1")

	coupons := v1.Group("/coupons")
	{
		coupons.POST("", app.Create)
		coupons.GET("", app.Get)
		coupons.POST("/basket", app.Apply)
	}

	return router
}

func (app *Application) Run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.Addr),
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("Listening",
		zap.String("on", srv.Addr),
	)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		app.logger.Errorw("start http server failed", "error", err)
		return err
	}

	return nil
}

func (app *Application) writeJSONResponse(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}

func (app *Application) writeJSONError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}
