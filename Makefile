.PHONY: dev build deploy

dev:
	GOEXPERIMENT=jsonv2 go run cmd/main.go
