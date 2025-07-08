# stage 1, build go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

ARG service_dir


COPY ${service_dir}/go.* ./
RUN go mod download

COPY ${service_dir}/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server

# Stage 2, final, small image
FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080 8081

ENTRYPOINT ["./server"]
