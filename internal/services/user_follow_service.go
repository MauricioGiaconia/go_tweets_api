package services

import (
	"database/sql"
	"fmt"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/repositories"
)

type FollowService struct {
	DB *sql.DB
}

func NewFollowService(db *sql.DB) *FollowService {
	return &FollowService{DB: db}
}

func (us *FollowService) FollowUser(follow *models.UserFollow) (bool, error) {
	userFollow, err := repositories.FollowUser(us.DB, follow)

	if err != nil {
		return userFollow, fmt.Errorf("Error followed user: %v", err)
	}

	return userFollow, nil
}

func (us *FollowService) GetFollows(userId *int64, relationType *string) (models.UserFollows, error) {
	userFollowers, err := repositories.GetFollows(us.DB, *userId, *relationType)

	if err != nil {
		return models.UserFollows{}, fmt.Errorf("Error getting followers: %v", err)
	}

	return *userFollowers, nil
}
