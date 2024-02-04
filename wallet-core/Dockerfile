FROM golang:1.21

WORKDIR /app/

RUN apt-get update && apt-get install -y librdkafka-dev

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./ ./

WORKDIR /app/cmd/wallet-core/
RUN go build -o web-server

CMD ["./web-server"]