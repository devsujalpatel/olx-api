.PHONY: build run
	
build:
	@go build -trimpath -ldflags="-s -w" -o bin/api ./cmd/api 

run: build
	@./bin/api

migrate-up: 
	@go run ./cmd/migrate up
	
migrate-down: 
	@go run ./cmd/migrate down