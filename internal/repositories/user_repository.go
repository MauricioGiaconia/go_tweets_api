package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
)

// CreateUser crea un nuevo usuario en la base de datos.
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

// GetUserById obtiene un usuario por su ID desde la base de datos.
func GetUserById(db *sql.DB, id int64) (models.User, error) {
	var user models.User
	err := db.QueryRow(`SELECT id, name, email, created_at FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, fmt.Errorf("[x] Error to get user: %v", err)
	}

	return user, nil
}
