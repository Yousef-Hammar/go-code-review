package main

import (
	"log"
	"os"
	"runtime"

	"go.uber.org/zap"

	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/api"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/repository/memory"
	"github.com/Yousef-Hammar/go-code-review/coupon_service/internal/service"
)

const (
	localEnv = "local"
	numCPU   = 32
)

func init() {
	if os.Getenv("env") != localEnv && runtime.NumCPU() != numCPU {
		log.Print("this api is meant to be run on 32 core machines")
		os.Exit(1)
	}
}

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
