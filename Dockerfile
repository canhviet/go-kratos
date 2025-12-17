FROM golang:1.25.5 as builder

WORKDIR /app
COPY . .
RUN go build -o myapp ./cmd/myapp

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/myapp .
COPY configs ./configs

EXPOSE 8000
CMD ["./myapp", "-conf", "./configs/config.yaml"]

