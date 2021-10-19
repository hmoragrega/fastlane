FROM golang:latest AS builder

ARG GOOS="linux"
ARG GOARCH="amd64"
ARG GOARM=""

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags="-w -s" -o /fastlane cmd/server/server.go

FROM alpine AS release

WORKDIR /

COPY --from=builder /fastlane /

CMD ["/fastlane"]
