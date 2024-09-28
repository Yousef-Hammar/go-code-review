package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memory"
	"coupon_service/internal/service"
)

const (
	localEnv = "local"
	numCPU   = 32
)

var (
	cfg  = config.New()
	repo = memory.New()
)

func init() {
	if os.Getenv("env") != localEnv && runtime.NumCPU() != numCPU {
		log.Print("this api is meant to be run on 32 core machines")
		os.Exit(1)
	}
}

func main() {
	svc := service.New(repo)
	本 := api.New(cfg.API, svc)
	本.Start()
	fmt.Println("Starting Coupon service server")
	<-time.After(1 * time.Hour * 24 * 365)
	fmt.Println("Coupon service server alive for a year, closing")
	本.Close()
}
