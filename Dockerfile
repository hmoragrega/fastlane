# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:alpine3.14 as build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM" > /log

RUN apk update && apk --no-cache add ca-certificates

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/fastlane"]
COPY fastlane /