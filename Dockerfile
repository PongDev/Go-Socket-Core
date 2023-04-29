FROM golang:latest AS base

WORKDIR /app

FROM base AS builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM base AS deploy
COPY --from=builder /app/main ./

CMD ["./main"]
