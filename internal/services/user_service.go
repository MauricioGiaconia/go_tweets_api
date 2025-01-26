package services

import (
	"database/sql"
	"fmt"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/repositories"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

func (us *UserService) CreateUser(user *models.User) (int64, error) {
	userID, err := repositories.CreateUser(us.DB, user)
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}
	return userID, nil
}

func (us *UserService) GetUserById(id int64) (models.User, error) {
	user, err := repositories.GetUserById(us.DB, id)

	if err != nil {
		return models.User{}, fmt.Errorf("Error fetching user: %v", err)
	}

	return user, nil
}
