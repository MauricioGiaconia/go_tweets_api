package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
)

func CreateUser(db *sql.DB, user *models.User) (int64, error) {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		log.Printf("[x] Error to create user: %v", err)
		return 0, fmt.Errorf("Error to create user")
	}

	return id, nil
}
