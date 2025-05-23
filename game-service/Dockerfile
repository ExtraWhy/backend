FROM golang:1.24 AS builder

WORKDIR /app

ADD bitvector/ bitvector
ADD cleopatra/ cleopatra
ADD slot-game/ slot-game
ADD slots/ slots
ADD util/ util

COPY go.mod go.sum main.go ./

RUN go mod download

RUN CGO_ENABLED=0 go build -o game-service -ldflags="-X 'main.version=0.0.1'" -tags netgo -a -v main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/game-service ./

CMD ["./game-service"]
