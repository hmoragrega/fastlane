SHELL=bash

#PLATFORM=linux/arm/v7,linux/amd64
PLATFORM=linux/arm/v7

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

server:
	@go run ./cmd/server

image:
	@docker buildx build \
	--push \
	--platform $(PLATFORM) \
	--tag hmoragrega/fastlane:latest \
	--tag hmoragrega/fastlane:0.0.1 .

vet:
	@go vet ./...

fmt:
	@go fmt ./...