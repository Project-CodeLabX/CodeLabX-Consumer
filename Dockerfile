FROM golang:1.21-alpine3.18 AS Builder

WORKDIR /app

copy . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd

EXPOSE 9010

FROM python:3.9.19-alpine 

WORKDIR /

COPY --from=Builder /api /api

EXPOSE 9010

CMD [ "/api" ]