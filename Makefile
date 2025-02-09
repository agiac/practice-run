.PHONY: start test generate

start:
	go run .

test:
	go test -race -v ./...

generate:
	go generate ./...