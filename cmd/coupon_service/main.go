package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/api"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/config"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/repository/memory"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/service"
)

func main() {
	cfg := config.New()
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	repo := memory.New()
	svc := service.New(repo)

	app := api.New(cfg, logger, svc)

	router := app.Mount(gin.DebugMode)

	if err := app.Run(router); err != nil {
		log.Fatal(err)
	}
}
