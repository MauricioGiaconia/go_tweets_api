package services

import (
	"database/sql"
	"fmt"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

func (us *UserService) CreateUser(user *models.User) (int64, error) {
	result, err := us.DB.Exec(`INSERT INTO users (name, email, password) VALUES (?, ?, ?)`,
		user.Name, user.Email, user.Password)
	if err != nil {
		return 0, fmt.Errorf("[x] Error inserting user: %v", err)
	}

	// Obtenemos el ID del último usuario insertado
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("[x] Error retrieving last inserted user ID: %v", err)
	}

	// Retornamos el ID del usuario y ningún error
	return userID, nil
}

func (us *UserService) GetUserById(id int64) (models.User, error) {
	var user models.User

	err := us.DB.QueryRow(`SELECT id, name, email, created_at FROM users WHERE id = ?`, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// Si no se encuentra el usuario, retornamos un error específico
			return models.User{}, fmt.Errorf("user with id %d not found", id)
		}
		return models.User{}, fmt.Errorf("[x] Error to get user: %v", err)
	}

	return user, nil
}
