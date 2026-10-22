# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Final stage
FROM alpine:3.20

WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /server /app/server
RUN mkdir -p /app/uploads

EXPOSE 8080
ENV PORT=8080
ENV UPLOAD_DIR=/app/uploads

CMD ["/app/server"]
