.PHONY: swag run build

swag:
	swag init -d cmd/kasir-api,internal/handler,internal/models,internal/service,internal/utils -g main.go --parseInternal

run: swag
	go run cmd/kasir-api/main.go

build: swag
	go build -o kasir-api cmd/kasir-api/main.go
