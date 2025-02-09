# UALA_BACKEND_CHALLENGE - Prueba técnica para Uala

## Descripción

Este proyecto es un challenge backend para Uala. La prueba consiste en crear una plataforma similar a *Twitter* donde se permitirá:

- Crear usuarios.
- Que los usuarios puedan seguirse entre sí.
- Postear tweets con un máximo de 280 caracteres.
- Que un usuario obtenga el timeline de todos los usuarios a los que sigue (es decir, obtener todos los tweets).

### Pasos previos

Esta API utiliza **SQLite** como base de datos SQL, la cual se ejecuta en memoria. También tiene soporte para **PostgreSQL**, que sería la base de datos principal en un entorno de *producción*.

Si deseas utilizar **PostgreSQL**, es necesario tener un archivo `.env` en la raíz del proyecto con las siguientes variables:

```bash
DB_HOST=...
DB_USER=...
DB_PASSWORD=...
DB_NAME=...
DB_SSLMODE=...
DB_PORT=...
```

**Nota:** _Si queres utilizar postgres, tenes que levantar el proyecto sin Docker._

## Ejecución del proyecto

Para ejecutar este proyecto, puedes elegir entre dos opciones. En ambas, debes ejecutar los comandos desde la terminal, ubicada en la raíz del proyecto:

- Utilizar Docker.
- Ejecutarlo de forma manual (sin Docker).

### Ejecución con Docker

La ejecución con Docker tiene la ventaja de configurar automáticamente Redis sin tener que hacer pasos adicionales.
Pasos para ejecutar con Docker:

1- Construir la imagen de Docker:

```bash
docker build -t uala_backend_challenge .   
```

2- Ejecuta un contenedor:
Para levantar el contenedor con Docker, usa el siguiente comando. Si deseas cargar variables de entorno desde un archivo .env, puedes agregar la instrucción --env-file .env.

**Nota:** La opción -d ejecuta el contenedor en segundo plano. Si deseas ver la consola de la API, quita la opción -d.

```bash
docker run --name uala-challenge -d -p 8080:8080 uala_backend_challenge
```

3- Probar la API:
Para verificar el correcto funcionamiento del proyecto, realiza una solicitud GET a este endpoint:

```bash
localhost:8080/ping
```

Si todo está bien, deberías recibir la siguiente respuesta:

```json
{
    "code": 200,
    "data": "Pong"
}
```

4- Detener el contenedor:
Para detener el contenedor, ejecuta el siguiente comando:

```bash
docker stop uala-challenge
```

5- Una vez detenido el contenedor, puedes eliminarlo con:

```bash
docker rm uala-challenge
```

---

### Ejecución SIN Docker

**Aclaraciones sobre Redis:**

- Si deseas usar Redis (para almacenar en caché la información consultada de manera recurrente), puedes levantar una instancia de Redis localmente.
- Si no tienes Redis instalado o no deseas usarlo, el proyecto funcionará igual, pero sin la funcionalidad de caché.

1- Paso 1: Descargar imagen Redis.

```bash
docker pull redis
```

2- Paso 2: Iniciar el contenedor en segundo plano en el puerto 6379 a partir de la imagen Redis obtenida.

```bash
docker run --name redis -p 6379:6379 -d redis
```

3- Paso 3 (opcional): Ejecuta dentro del contenedor redis, la herramienta de comandos redis-cli, esto te permitirá interactura con el servidor Redis desde tu terminal en caso de que lo requieras.

```bash
docker exec -it redis redis-cli
```

#### Instalación de dependencias

```bash
go mod tidy
```

Ejecución de la API.

**Nota**: Si no envías la variable de entorno **--db**, la API usará SQLite por defecto. Si deseas usar PostgreSQL, agrega --db=postgres al comando de ejecución. Tambien, si no envias la variable **--port**, se tomará el valor 8080 por defecto.

```bash
go run cmd/api/main.go --db=sqlite --port=8080
```

Deberías ver el siguiente mensaje indicando que la API está en funcionamiento:

```bash
[GIN-debug] Listening and serving HTTP on :8080
```

Para verificar que la API está funcionando correctamente, realiza una solicitud GET al siguiente endpoint:

```bash
localhost:8080/ping
```

Si la API esta funcionando correctamente, deberías ver el siguiente mensaje como respuesta:

```json
{
    "code":200,
    "data":"pong"
}
```
