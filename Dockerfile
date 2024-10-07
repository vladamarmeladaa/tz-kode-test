FROM golang:1.21.6-alpine

RUN apk add --no-cache git

WORKDIR /app/tz-kode

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/tz-kode ./cmd/main.go

EXPOSE 8080

CMD ["./out/tz-kode"]