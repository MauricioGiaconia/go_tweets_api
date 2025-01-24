package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDatabase define la conexión para SQLite
type SQLiteDatabase struct{}

func (s *SQLiteDatabase) Connect() (*sql.DB, error) {
	// Usar SQLite en memoria (se perderá al cerrar la aplicación)
	connStr := "file::memory:?cache=shared" // Esto crea una base de datos en memoria
	// Intentar abrir la conexión a la base de datos SQLite en memoria
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("[x] DB error to open SQLite in memory: %v", err)
	}

	ConfigurePoolConnection(db)

	// Verificar si la conexión está activa
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("[x] SQLite DB ping error: %v", err)
	}

	// Una vez validada la conexión, creo las tablas necesarias para que el proyecto funcione
	err = createTables(db)

	if err != nil {
		return nil, fmt.Errorf("[x] SQLite Error creating tables: %v", err)
	}

	fmt.Printf("[x] SQLiteDB connection success")
	// En este punto, la conexión está lista para usarse.
	return db, nil
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
			followed_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY(follower_id, followed_id),
			FOREIGN KEY(follower_id) REFERENCES users(id),
			FOREIGN KEY(followed_id) REFERENCES users(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("[x] Error creating follows table: %v", err)
	}

	return nil
}
