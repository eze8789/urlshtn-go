FROM golang:1.15-alpine as builder
LABEL maintainer="Ezequiel Reyna <ezequiel.reyna@gmail.com>"

WORKDIR /app

COPY  go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o urlshtn-go -ldflags="-s -w" .

# Stage for the application
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/urlshtn-go /app/configs/app/config.yaml /app/configs/sql/create_url_shortener.sql /app/

EXPOSE 8080

ENV CONFIG=/app/config.yaml

CMD /app/urlshtn-go -config ${CONFIG}