# syntax=docker/dockerfile:experimental
FROM arm32v7/alpine:3.11 as build

RUN apk update && apk --no-cache add ca-certificates

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/fastlane"]
COPY fastlane /