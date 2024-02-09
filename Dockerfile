FROM golang:1.21.6-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
ENV GOPROXY https://goproxy.cn,direct
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./main ./...



FROM alpine:latest AS runner
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app .
EXPOSE 8081
CMD ["./main"]
