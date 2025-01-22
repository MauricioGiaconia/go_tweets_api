package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func Init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Configurar el pool de conexiones
	DB.SetMaxOpenConns(25)                 // Máximo número de conexiones abiertas
	DB.SetMaxIdleConns(25)                 // Máximo número de conexiones inactivas
	DB.SetConnMaxLifetime(5 * time.Minute) // Tiempo máximo de vida de una conexión

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MySQL DB: Successful connection!")
}
