#!make
include ./.env
export $(shell sed 's/=.*//' ./.env)

run-unit-tests:
	go test -v ./...

run-integration-tests:
	LONG=true go test -v ./...

generate-mocks:
	mockery

docker-build:
	 docker build -t coupon-service .

docker-run: docker-build
	docker run --env-file .env -p ${ADDR}:${ADDR} coupon-service

