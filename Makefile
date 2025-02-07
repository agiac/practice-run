start:
	go run .

test:
	go test -race ./...

generate:
	go generate ./...