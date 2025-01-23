package services

import (
	"database/sql"
	"fmt"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/repositories"
)

// CreateUserService valida los datos y crea un nuevo usuario
func CreateUserService(db *sql.DB, user *models.User) (int64, error) {
	// Validación de campos
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return 0, fmt.Errorf("All fields are required")
	}

	userID, err := repositories.CreateUser(db, user)
	if err != nil {
		return 0, fmt.Errorf("Error in CreateUserService: %v", err)
	}

	return userID, nil
}
