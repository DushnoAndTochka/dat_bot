FROM golang:1.18-alpine

WORKDIR /app
RUN apk add --no-cache bash

COPY go.mod ./
COPY go.sum ./
COPY .env ./.env
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

COPY Taskfile.yml ./

# RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go build -o /dushno_and_tochka_bot ./cmd/dushno_and_tochka_bot/main.go


CMD [ "/dushno_and_tochka_bot" ]