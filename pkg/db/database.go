package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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
