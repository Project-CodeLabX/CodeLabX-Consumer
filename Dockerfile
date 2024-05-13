FROM golang:1.21-alpine3.18 AS Builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd

FROM  openjdk:17-alpine

WORKDIR /

COPY --from=Builder /api /api

EXPOSE 9020

CMD ["/api"]