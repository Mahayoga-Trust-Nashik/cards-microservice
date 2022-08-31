FROM golang:latest as builder

WORKDIR /app
COPY . .

RUN go mod download

RUN env CGO_ENABLED=0 go build -o /cards-microservice

FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ=Asia/Kolkata

WORKDIR /

COPY --from=builder /cards-microservice /cards-microservice

ENV MYSQL_CONNECTION="mahayoga_api:h7*Z537pdHjA@tcp(mahayoga-database:3306)/mahayoga_mobile?charset=utf8mb4&parseTime=True&loc=Local"

EXPOSE 5300

ENTRYPOINT ["/cards-microservice"]