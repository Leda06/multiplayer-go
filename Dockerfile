FROM golang:1.23.2-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
ARG GOPROXY
ARG GO111MODULE
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./cmd/entry ./...



FROM alpine:latest AS runner
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app .
EXPOSE 8081
CMD ["./cmd/entry/entry"]
