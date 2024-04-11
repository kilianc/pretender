test:
	go test -v ./internal/...

build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/pretender cmd/main.go

run:
	go run cmd/main.go --responses README.md

docker-build:
	docker build . -t pretender:latest

docker-run:
	docker run --rm -v $(PWD)/README.md:/README.md -p 8080:8080 pretender:latest --responses /README.md

docker: docker-build docker-run

.PHONY: test build run docker-build docker-run docker
