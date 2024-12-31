FROM golang:1.21-alpine

WORKDIR /app

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/ .

RUN go build -o main .

ENV GIN_MODE=debug

EXPOSE 8081

CMD ["./main"] 