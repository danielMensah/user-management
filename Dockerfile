# User
FROM alpine:3.13.1 as user
ARG uid=10001
ARG gid=10001
RUN echo "scratchuser:x:${uid}:${gid}::/home/scratchuser:/bin/sh" > /scratchpasswd

# Certs
FROM alpine:3.13.1 as certs
RUN apk add -U --no-cache ca-certificates

# Build
FROM golang:1.18-alpine as build
ENV CGO_ENABLED 0
ENV GO111MODULE=on

WORKDIR /code/
COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY ./go.mod/ ./go.sum/ ./
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o ./bin/api ./cmd/api

# Entrypoints
FROM scratch as api
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=user /scratchpasswd /etc/passwd
COPY --from=build /code/bin/api .
USER scratchuser
EXPOSE 8000
ENTRYPOINT ["/api"]
