FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod .
#COPY go.sum .
RUN go mod download

COPY . .
COPY config/yml/config.yml /app/config/config.yml

RUN go build -o /app/smartcompost-webapp ./cmd

EXPOSE 8080

CMD ["/app/smartcompost-webapp"]
