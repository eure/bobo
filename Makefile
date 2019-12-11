.PHONY: init build build-arm6 dev lint test

GO111MODULE=on

init:
	go mod download

# build binary
build:
	go build -o bin/bobo ./cmd/examplebot


# build binary for Raspberry Pi
build-arm6:
	GOOS=linux GOARCH=arm GOARM=6 go build -o bin/arm6 ./cmd/examplebot

# run bot on local environment
dev:
	go run ./cmd/examplebot

# Exec golint, vet, gofmt
lint:
	@type golangci-lint > /dev/null || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

test:
	@type gosec > /dev/null || go get github.com/securego/gosec/cmd/gosec
	gosec -quiet ./...
	go test ./... -count=1;
