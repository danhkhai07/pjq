# Dockerfile for pjq job queue service
FROM golang:1.25.6-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/pjq-app ./cmd/main.go

CMD ["pjq-app"]
