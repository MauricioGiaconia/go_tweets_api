package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// PostgresDatabase define la conexión para PostgreSQL
type PostgresDatabase struct{}

// Para simplificar la prueba técnica, se utilizará una base de datos SQLite en memoria,
// pero el código está preparado para conectarse a PostgreSQL como base de datos principal.
// Para la prueba técnica, se usa SQLite en memoria (modo ':memory:') para simular la DB PostgreSQL.
// Sin embargo, en un escenario de producción o para un sistema escalable, PostgreSQL es una buena opción
// por llo siguiente:
//   - **Escalabilidad**: PostgreSQL puede manejar grandes volúmenes de datos y muchos usuarios concurrentes,
//     lo que lo hace ideal para aplicaciones a gran escala como la que se propone en esta prueba.
//   - **Consistencia y durabilidad**: Con PostgreSQL se garantiza la consistencia de los datos y una mayor
//     fiabilidad al manejar operaciones complejas, lo que es crucial en aplicaciones donde la integridad de los datos
//     es fundamental.
//   - **Rendimiento optimizado**: PostgreSQL tiene un motor de consultas avanzado y soporta índices, optimizaciones
//     de consultas y otras características que lo hacen adecuado para aplicaciones que requieren un alto rendimiento.

func (p *PostgresDatabase) Connect() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("[x] Cannot load .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	if dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" || dbHost == "" {
		return nil, fmt.Errorf("[x] PostgreSQL DB missing variables to connect")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	// Intenta abrir la conexión a la base de datos
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("[x] PostgreSQL DB error to open: %v", err)
	}

	ConfigurePoolConnection(db)
	// Se verifica la conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("[x] PostgreSQL DB ping error: %v", err)
	}

	return db, nil
}
