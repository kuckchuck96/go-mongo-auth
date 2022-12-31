dev:
	go run --race cmd/go-mongo-auth/main.go --profile=dev

lint:
	golangci-lint run

swag:
	swag init -g ./cmd/go-mongo-auth/main.go