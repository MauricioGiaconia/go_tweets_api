# Benchmark de rendimiento de la obtención de un timeline

## Requisitos previos

1. **Base de datos PostgreSQL**:
   - Asegúrate de tener PostgreSQL corriendo y accesible en tu entorno.
   - El archivo `.env` debe contener las configuraciones correctas para conectar a la base de datos PostgreSQL y Redis.

2. **Archivo `.env`**:
   - Debes tener un archivo `.env` en la carpeta *benchmark* del proyecto con las siguientes variables de entorno configuradas (reemplazar los valores con lo que corresponda en tu caso):

3. **Relación de seguidores**:
  - Para realizar correctamente el benchmark, es necesario que los usuarios en la base de datos tengan una relación de seguidores configurada.
    **El usuario con ID 3 debe estar siguiendo al usuario con ID 1.**  
    Esto es importante para realizar las pruebas de consulta de timelines de un usuario específico.

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=1234
DB_NAME=uala_backend_challenge
DB_SSLMODE=disable
DB_PORT=9090

Pasos para utilizar Redis.

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

Finalmente, podes detener Redis para hacer un benchmark sin esa tecnología y ver la diferencia que existe entre usar unicamente una DB sql vs combinar el uso SQL con Redis.

```bash
docker stop redis
```
