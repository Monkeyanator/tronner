FROM golang:1.15.6 AS builder
WORKDIR /app
COPY . .
RUN make server.build
RUN make wasm.build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY . /app/
COPY --from=builder /app/server /app/
EXPOSE 8080
WORKDIR /app/
CMD ["./server", "--port", "8080"]
