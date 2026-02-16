# Prerequisites:
# - Linux: libasound2-dev (ALSA)
# - macOS: none (CoreAudio is built-in)
# - Windows: none (WASAPI is built-in)

VERSION ?= dev
BINARY_NAME = speak
LDFLAGS = -ldflags "-s -w -X main.Version=$(VERSION)"

.PHONY: build clean test

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/speak/

test:
	go test ./... -v -timeout 30s

clean:
	rm -f $(BINARY_NAME)
	rm -f speak-*

# Native build for current platform (used by CI)
build-native:
	go build $(LDFLAGS) -o $(BINARY_NAME)-$(shell go env GOOS)-$(shell go env GOARCH)$(shell go env GOEXE) ./cmd/speak/
