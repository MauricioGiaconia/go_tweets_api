package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	Connect() (*sql.DB, error) // Connect realiza la conexión a la base de datos.
}

// CloseDatabase cierra la conexión a la base de datos
func CloseDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("[x] DB error to close connection: %v", err)
	}
	log.Println("[x] DB Connection closed.")
}

// Configuración del pool de conexiones
func ConfigurePoolConnection(db *sql.DB) {
	db.SetMaxOpenConns(10)                  // Número máximo de conexiones abiertas
	db.SetMaxIdleConns(5)                   // Número máximo de conexiones inactivas
	db.SetConnMaxLifetime(30 * time.Minute) // Tiempo máximo de vida de la conexión
}
