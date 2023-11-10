run:
	@go run main.go
build:
	@go build

.PHONY: all build deps lint test
