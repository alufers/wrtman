############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git nodejs yarn bash && mkdir -p /build/wrtman

WORKDIR /build/wrtman

COPY go.mod go.mod
COPY go.sum go.sum
COPY ./wrtman-frontend/package.json ./wrtman-frontend/package.json
COPY ./wrtman-frontend/yarn.lock ./wrtman-frontend/yarn.lock

RUN go mod download -json && cd wrtman-frontend 

COPY . .

RUN mkdir -p /app && go generate && go build -tags embed_frontend -ldflags="-s -w" -o /app/wrtman

FROM scratch AS bin-unix
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/wrtman /app/wrtman
COPY --from=builder /build/wrtman/data/oui.txt /app/data/oui.txt

ENTRYPOINT ["/app/wrtman"]
