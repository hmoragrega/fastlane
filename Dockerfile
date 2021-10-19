FROM golang:latest AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /fastlane cmd/server/server.go

FROM alpine AS release

WORKDIR /

COPY --from=builder /fastlane /

CMD ["/fastlane"]
