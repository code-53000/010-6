FROM golang:1.21-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
ENV GONOSUMDB=*
ENV GOFLAGS=-mod=mod

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /seed-queue-core ./cmd/main.go

FROM alpine:latest

WORKDIR /

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

COPY --from=builder /seed-queue-core /seed-queue-core

EXPOSE 8080

CMD ["/seed-queue-core"]
