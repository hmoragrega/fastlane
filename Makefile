SHELL=bash

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

server:
	@go run ./cmd/server

vet:
	@go vet ./...

fmt:
	@go fmt ./...