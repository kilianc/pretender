build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/pretender cmd/pretender/main.go

test:
	go test -v ./internal/...

run:
	go run cmd/main.go --responses README.md

docker-build:
	docker build . -t pretender:latest

docker-run:
	docker run --rm -v $(PWD)/README.md:/README.md -p 8080:8080 pretender:latest --responses /README.md

docker:
	docker-build docker-run

version-check:
	go run scripts/versioncheck.go

.PHONY: build docker docker-build docker-run run test version-check
