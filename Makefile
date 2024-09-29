run-tests:
	go test ./...

generate-mocks:
	mockery

docker-build:
	 docker build -t coupon-service .

docker-run:
	docker run -p 8080:8080 coupon-service