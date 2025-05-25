FROM golang:1.24 AS builder

WORKDIR /app

ADD handlers/ handlers
COPY go.mod go.sum user-service.yaml main.go ./

RUN go mod download

RUN CGO_ENABLED=0 go build -o user-service -ldflags="-X 'main.version=0.0.1'" -tags netgo -a -v main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/user-service /app/user-service.yaml ./

CMD ["./user-service", "user-service.yaml"]
