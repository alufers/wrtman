############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git nodejs yarn && mkdir -p /build/wrtman

WORKDIR /build/wrtman

COPY go.mod go.mod
COPY go.sum go.sum
COPY ./wrtman-frontend/package.json ./wrtman-frontend/package.json
COPY ./wrtman-frontend/yarn.lock ./wrtman-frontend/yarn.lock

RUN go mod download -json

COPY . .

RUN mkdir -p /app/wrtman && go generate && go build -tags embed_frontend -ldflags="-s -w" 
