FROM golang:1.21.3

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /delete-service

EXPOSE 8124

CMD ["/delete-service"]