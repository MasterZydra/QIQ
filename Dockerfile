# Build stage
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -v -o . ./...

# Runtime stage
FROM alpine:latest
WORKDIR /app

ENV PORT=8080
ENV DOC_ROOT=/var/www/html
ENV DEV=false

RUN mkdir -p /var/www/html

COPY --from=builder /app/qiq .
COPY --from=builder /app/docker/ /var/www/html/
COPY --from=builder /app/doc/Rabbit.svg /var/www/html/

CMD ["sh", "-c", "./qiq -S 0.0.0.0:${PORT} $( [ \"$DEV\" = \"true\" ] && echo \"-dev\" ) -t ${DOC_ROOT}"]
