# Etapa de construcción
FROM golang:1.23.4-alpine3.21 as builder

RUN apk add --no-cache --no-cache \
    git \
    build-base \
    linux-headers \
    gcc \
    musl-dev

# Establecer el directorio de trabajo en el contenedor
WORKDIR /app

# Copiar go.mod y go.sum primero (para aprovechar el cache de Docker)
COPY go.mod go.sum ./ 

# Descargar las dependencias
RUN go mod tidy

# Copiar el resto de los archivos fuente (manteniendo la estructura de directorios)
COPY . . 

# Compilar el binario
RUN go build -o /bin/api cmd/api/main.go

# Etapa final - imagen ligera
FROM alpine:latest

# Instalar Redis y dependencias necesarias para la API
RUN apk add --no-cache \
    redis \
    bash \
    && mkdir /app

# Copiar el binario compilado desde la etapa de construcción
COPY --from=builder /bin/api /bin/api

# Copiar archivos de la aplicación
COPY . /app

# Exponer los puertos necesarios
EXPOSE 8080
EXPOSE 6379

# Crear un script de inicio para ejecutar Redis y la API
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Comando para ejecutar el script de inicio
CMD ["/start.sh"]
