
FROM golang:1.19.1-alpine3.16 as builder

ENV GOOS linux
ENV CGO_ENABLED 0

WORKDIR /app

ARG DB_PORT
ARG SECRET
ARG PORT
ARG GIN_MODE

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd

FROM alpine:3.16.2 as production

RUN apk add --no-cache ca-certificates

COPY --from=builder app .

EXPOSE 8080



CMD ./app