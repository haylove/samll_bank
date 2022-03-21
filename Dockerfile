FROM golang:1.17.4-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go build -ldflags "-s -w" -o main

FROM alpine:3.15
WORKDIR /app

COPY --from=builder /app/main .
COPY /app.env .

EXPOSE 8080
CMD [ "/app/main" ]