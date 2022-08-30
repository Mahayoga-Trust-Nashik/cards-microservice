FROM golang:latest as builder

WORKDIR /app
COPY . .

RUN go mod download

RUN env CGO_ENABLED=0 go build -o /cards-microservice

FROM alpine:latest

WORKDIR /

COPY --from=builder /cards-microservice /cards-microservice

EXPOSE 5300

ENV TZ="Asia/Kolkata"

ENTRYPOINT ["/cards-microservice"]