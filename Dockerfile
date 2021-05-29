FROM golang:1.15.6 AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY . /app/
COPY --from=builder /app/main /app/
EXPOSE 8080
WORKDIR /app/
CMD ["./main", "--port", "8080"]