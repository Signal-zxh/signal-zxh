FROM golang:1.26.4-alpine3.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o signalzxh-blog .

FROM alpine:3.24 AS final

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/signalzxh-blog .

EXPOSE 8080
CMD ["./signalzxh-blog"]
