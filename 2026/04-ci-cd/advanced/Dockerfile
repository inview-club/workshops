FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/products ./cmd/products

FROM alpine:3.18

RUN adduser -D -u 10001 appuser
USER appuser

COPY --from=builder /bin/products /products

EXPOSE 30001

ENTRYPOINT ["/products"]
