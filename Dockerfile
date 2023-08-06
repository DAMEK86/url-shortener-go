# syntax=docker/dockerfile:1

FROM golang:1.20 AS builder
# smoke test to verify if golang is available
RUN go version

ARG BUILD_VERSION=local
ARG SOURCE_VERSION=local
ARG CGO_ENABLED=0

COPY . /build-dir/
WORKDIR /build-dir/
RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-w -s \
    -X github.com/damek86/url-shortener-go/internal/config.BuildVersion=${BUILD_VERSION} \
    -X github.com/damek86/url-shortener-go/internal/config.SourceVersion=${SOURCE_VERSION}" \
    -a -o app internal/main.go
RUN go test -cover -v ./...

FROM alpine:3

ARG UID=1001
ARG USER=app
ARG GID=1001
ARG GROUP=app
ENV WORKINGDIR /app

EXPOSE 8080

RUN apk --no-cache add ca-certificates

WORKDIR $WORKINGDIR
RUN addgroup -g $GID -S $GROUP && adduser -u $UID -S $USER -G $GROUP && \
    mkdir -p /app &&\
    chown -R $USER:$GROUP /app

USER $USER
COPY --from=builder /build-dir/staticcontent/swagger-ui ./staticcontent/swagger-ui
COPY --from=builder /build-dir/app url-shortener

CMD ./url-shortener