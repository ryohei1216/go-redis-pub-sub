## build stage ##
FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o build-file /app/cmd/main.go

## run stage ##
FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app .

EXPOSE 8080

CMD [ "/app/build-file" ]