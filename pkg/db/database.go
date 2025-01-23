package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// SetupSQLiteDatabase configura la conexión a la base de datos SQLite en memoria.
func SetupSQLiteDatabase() (*sql.DB, error) {
	// Usar SQLite en memoria (se perderá al cerrar la aplicación)
	connStr := "file::memory:?cache=shared" // Esto crea una base de datos en memoria
	// Intentar abrir la conexión a la base de datos SQLite en memoria
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("[x] DB error to open SQLite in memory: %v", err)
	}

	// Configuración del pool de conexiones
	db.SetMaxOpenConns(10)                  // Número máximo de conexiones abiertas
	db.SetMaxIdleConns(5)                   // Número máximo de conexiones inactivas
	db.SetConnMaxLifetime(30 * time.Minute) // Tiempo máximo de vida de la conexión

	// Verificar si la conexión está activa
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("[x] DB ping error: %v", err)
	}

	// Una vez validada la conexión, creo las tablas necesarias para que el proyecto funcione
	err = createTables(db)

	if err != nil {
		return nil, fmt.Errorf("[x] Error creating tables: %v", err)
	}

	fmt.Printf("[x] SQLiteDB connection success")
	// En este punto, la conexión está lista para usarse.
	return db, nil
}

// SetupDatabase configura la conexión a la base de datos.
// Para simplificar la prueba técnica, se utilizará una base de datos SQLite en memoria,
// pero el código está preparado para conectarse a PostgreSQL como base de datos principal.
// Para la prueba técnica, se usa SQLite en memoria (modo ':memory:') para simular la DB PostgreSQL.
// Sin embargo, en un escenario de producción o para un sistema escalable, PostgreSQL es una buena opción
// por llo siguiente:
// - **Escalabilidad**: PostgreSQL puede manejar grandes volúmenes de datos y muchos usuarios concurrentes,
//   lo que lo hace ideal para aplicaciones a gran escala como la que se propone en esta prueba.
// - **Consistencia y durabilidad**: Con PostgreSQL se garantiza la consistencia de los datos y una mayor
//   fiabilidad al manejar operaciones complejas, lo que es crucial en aplicaciones donde la integridad de los datos
//   es fundamental.
// - **Rendimiento optimizado**: PostgreSQL tiene un motor de consultas avanzado y soporta índices, optimizaciones
//   de consultas y otras características que lo hacen adecuado para aplicaciones que requieren un alto rendimiento.

func SetupDatabase() (*sql.DB, error) {
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
		return nil, fmt.Errorf("[x] DB missing variables to connect")
	}

	// Crear el string de conexión
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	// Intenta abrir la conexión a la base de datos
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("[x] DB error to open: %v", err)
	}
	// Configuración del pool de conexiones
	db.SetMaxOpenConns(10)                  // Número máximo de conexiones abiertas
	db.SetMaxIdleConns(5)                   // Número máximo de conexiones inactivas
	db.SetConnMaxLifetime(30 * time.Minute) // Tiempo máximo de vida de la conexión

	// Se verifica la conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("[x] DB ping error: %v", err)
	}

	fmt.Print("Hola bbto")
	return db, nil
}

// CloseDatabase cierra la conexión a la base de datos
func CloseDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("[x] DB error to close connection: %v", err)
	}
	log.Println("[x] DB Connection closed.")
}

// createTables crea las tablas necesarias en SQLite
// Nota: Esto solo sirve para levantar la DB en memoria. Si se utiliza PostgreSQL, se deben crear las tablas en su respectiva DB
func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("[x] Error creating users table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tweets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("[x] Error creating tweets table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS follows (
			follower_id INTEGER NOT NULL,
			following_id INTEGER NOT NULL,
			PRIMARY KEY(follower_id, following_id),
			FOREIGN KEY(follower_id) REFERENCES users(id),
			FOREIGN KEY(following_id) REFERENCES users(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("[x] Error creating follows table: %v", err)
	}

	return nil
}
