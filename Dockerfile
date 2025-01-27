# Etapa de construcción
FROM golang:1.23.4-alpine3.21 as builder

RUN apk add --no-cache --no-cache \
    git \
    build-base \
    linux-headers \
    gcc \
    musl-dev

WORKDIR /app

# Copiar go.mod y go.sum primero (para aprovechar el cache de Docker)
COPY go.mod go.sum ./ 

RUN go mod tidy

COPY . . 

RUN go build -o /bin/api cmd/api/main.go

FROM alpine:latest

# Instalar Redis y dependencias necesarias para la API
RUN apk add --no-cache \
    redis \
    bash \
    && mkdir /app

COPY --from=builder /bin/api /bin/api

COPY . /app

# Exponer los puertos necesarios
EXPOSE 8080
EXPOSE 6379

COPY start.sh /start.sh
RUN chmod +x /start.sh

CMD ["/bin/sh", "/start.sh"]