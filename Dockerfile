FROM golang:1.26.4-alpine3.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o signal-zxh .

FROM alpine:3.24 AS final


WORKDIR /app
COPY --from=builder /app/signal-zxh .

EXPOSE 8080
CMD ["./signal-zxh"]
