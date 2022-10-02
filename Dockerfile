# Build Stage
FROM golang:1.18.6-alpine3.16 AS builder

RUN apk update && \
    apk --no-cache add git && \
    apk --no-cache add curl

WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY app.env .
COPY *.sh .
COPY db/migration ./migration

EXPOSE 9000
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]