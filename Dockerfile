# -------- Stage 1 --------

FROM golang:1.24-alpine AS builder

WORKDIR /app

# caching go env
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o top10 .

# -------- Stage 2 --------

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/top10 .
COPY --from=builder /app/resource ./resource

EXPOSE 8080

CMD ["./top10"]