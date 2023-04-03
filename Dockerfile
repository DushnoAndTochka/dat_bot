FROM golang:1.20-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

COPY Taskfile.yml ./

# RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN env GOOS=linux GOARCH=amd64 go build -o /dushno_and_tochka_bot ./cmd/dushno_and_tochka_bot/main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /dushno_and_tochka_bot /dushno_and_tochka_bot

CMD [ "/dushno_and_tochka_bot" ]