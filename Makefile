.DEFAULT_GOAL := build
.PHONY: fmt vet build install clean test

APP_NAME := golox
INSTALL_DIR := ~/bin/

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o $(APP_NAME) ./cmd/golox/main.go

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(APP_NAME) $(INSTALL_DIR)

clean:
	go clean
	rm -f $(APP_NAME)

test:
	go test ./...
