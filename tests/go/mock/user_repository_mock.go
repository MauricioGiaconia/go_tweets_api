package userrepositorymock

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
)

type MockUserRepository struct{}

func (repo *MockUserRepository) GetUserById(db *sql.DB, id int64) (models.User, error) {
	if id == 1 {
		return models.User{
			ID:        1,
			Name:      "Mauricio Giaconia",
			Email:     "maurigiaconia@hotmail.com",
			CreatedAt: time.Now(),
		}, nil
	}
	return models.User{}, fmt.Errorf("User not found")
}
