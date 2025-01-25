package services

import (
	"database/sql"
	"fmt"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/repositories"
)

type TweetService struct {
	DB *sql.DB
}

func NewTweetService(db *sql.DB) *TweetService {
	return &TweetService{DB: db}
}

func (ts *TweetService) GetUserTimeline(followerId *int64, limit *int64, offset *int64) ([]models.Tweet, error) {

	timeline, err := repositories.GetTweetsTimeline(ts.DB, followerId, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("Error getting timeline: %v", err)
	}

	return timeline, nil
}

func (ts *TweetService) CountTimeline(followerId *int64) (int64, error) {
	total, err := repositories.CountTweetsTimeline(ts.DB, followerId)

	if err != nil {
		return 0, fmt.Errorf("Error counting timeline: %v", err)
	}

	return total, nil
}

func (ts *TweetService) PostTweet(tweet *models.Tweet) (bool, error) {
	const maxCharacters int = 280 //Maximos de caracteres permitidos en un tweet

	if len(tweet.Content) > maxCharacters {
		return false, fmt.Errorf("The content of the tweet must not exceed 280 characters")
	}

	tweetPosted, err := repositories.PostTweet(ts.DB, tweet)

	if err != nil {
		return false, fmt.Errorf("Error getting followers: %v", err)
	}

	return tweetPosted, nil
}

// Esta funcion, a diferencia del timeline, solo obtiene los tweets del usuario que los posteo (osea, los propios)
func (ts *TweetService) GetTweetsByUserId(userId *int64) ([]models.Tweet, error) {

	tweets, err := repositories.GetTweetsByUserId(ts.DB, userId)

	if err != nil {
		return nil, fmt.Errorf("Error getting user tweets: %v", err)
	}

	return tweets, nil
}
