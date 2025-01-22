# UALA_BACKEND_CHALLENGE - Prueba técnica para Uala

## Descripción

Proyecto que servirá como base para el desarrollo de futuras APIs con el objetivo de agilizar el desarrollo de la misma y facilitar el seguimiento del estandar API.

## Ejecutar el proyecto

La API usa variables de entornos para poder acceder a las distintas credenciales que se utilizarán, para poder conectarnos a la DB se deberá crear un archivo .env el cual debera contener la siguiente información:

```bash
DB_USER=example
DB_PASS=example
DB_HOST=localhost
DB_PORT=xxxx
DB_NAME=db_name
```

Para poder ejecutar la api y empezar a utilizarla, se deberá abrir un terminal en la raíz del proyecto y ejecutar los siguientes comandos:

Instalación de dependencias:

```bash
go mod tidy
```

```bash
go run cmd/api/main.go
```

En la consola de su terminal deberia aparecer el siguiente mensaje indicando que esta todo listo para empezar a probar los distintos endpoints:

```bash
[GIN-debug] Listening and serving HTTP on :8080
```

Para probar el correcto funcionamiento de la api, puede dirigirse a su navegador o a _postman_ y ejecutar el siguiente endpoint:

```bash
localhost:8080/ping
```

Si la API esta funcionando correctamente, debería ver el siguiente mensaje como respuesta:

```json
{
    "code":200,
    "data":"pong"
}
```