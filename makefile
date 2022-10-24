dev:
	go run cmd/go-mongo-auth/main.go --profile=dev

lint:
	golangci-lint run