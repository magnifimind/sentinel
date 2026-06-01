FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /sentinel-api ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /sentinel-api /usr/local/bin/sentinel-api
COPY migrations/ /app/migrations/

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["sentinel-api"]
