FROM golang:1.24 AS builder

WORKDIR /app

ADD models/ models
ADD player-cache/ player-cache
ADD proto-client/ proto-client
ADD rest-server/ rest-server
ADD serv-interface/ serv-interface
ADD ws-server/ ws-server

COPY go.mod go.sum requests-service.yaml main.go ./

RUN go mod download

RUN CGO_ENABLED=0 go build -o request-service -ldflags="-X 'main.version=0.0.1'" -tags netgo -a -v main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/request-service /app/requests-service.yaml ./

CMD ["./request-service", "requests-service.yaml"]
