server:
	go run cmd/main.go

production:
	go build ./cmd/main.go
	GIN_MODE=release ./main