.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build "-ldflags=-s -w -buildid=" -trimpath -o rakup ./cmd/rakup/main.go