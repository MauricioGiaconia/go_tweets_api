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

func (ufs *FollowService) FollowUser(follow *models.UserFollow) (bool, error) {
	userFollow, err := repositories.FollowUser(ufs.DB, follow)

	if err != nil {
		return userFollow, fmt.Errorf("Error followed user: %v", err)
	}

	return userFollow, nil
}

func (ufs *FollowService) GetFollows(userId *int64, relationType *string) (models.UserFollows, error) {

	// Validar que el relationType sea el adecuado segun la logica implementada en el repository
	if *relationType != "followers" && *relationType != "following" {

		return models.UserFollows{}, fmt.Errorf("Invalid follow type. Must be 'followers' or 'following'")
	}
	userFollowers, err := repositories.GetFollows(ufs.DB, *userId, *relationType)

	if err != nil {
		return models.UserFollows{}, fmt.Errorf("Error getting followers: %v", err)
	}

	return *userFollowers, nil
}
