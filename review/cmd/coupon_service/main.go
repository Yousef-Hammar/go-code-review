package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/api"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/repository/memory"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/service"
)

func main() {
	config := api.New()
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	repo := memory.New()
	svc := service.New(repo)

	app := api.NewApplication(config, logger, svc)

	router := app.Mount()

	if err := app.Run(router); err != nil {
		log.Fatal(err)
	}
}
