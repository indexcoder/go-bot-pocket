FROM golang:1.23.0-alpine3.19 AS builder

COPY . /telegram-bot-pocket/
WORKDIR /telegram-bot-pocket/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /telegram-bot-pocket/bin/bot .
COPY --from=0 /telegram-bot-pocket/config config/

EXPOSE 80

CMD ["./bot"]