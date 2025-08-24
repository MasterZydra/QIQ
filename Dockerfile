FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -v -o . ./...

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/qiq .
CMD ["./app-binary"]
