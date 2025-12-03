from golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main-api ./cmd/api/main.go
RUN go build -o main-worker ./cmd/worker/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main-api .
COPY --from=builder /app/main-worker .

EXPOSE 8080

CMD ["./main-api"]
