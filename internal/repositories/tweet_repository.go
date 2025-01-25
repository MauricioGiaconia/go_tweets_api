package repositories

import (
	"database/sql"
	"fmt"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
)

func PostTweet(db *sql.DB, tweet *models.Tweet) (bool, error) {
	query := `INSERT INTO tweets (user_id, content) VALUES ($1, $2)`

	_, err := db.Exec(query, tweet.UserID, tweet.Content)
	if err != nil {
		fmt.Println("[x] Error to create follow: %v", err)
		return false, fmt.Errorf("[x] Error to create follow: %v", err)
	}

	return true, nil
}

func GetTweetsByUserId(db *sql.DB, userId *int64) ([]models.Tweet, error) {
	query := `SELECT * FROM tweets WHERE user_id = $1`

	rows, err := db.Query(query, userId)

	if err != nil {
		return nil, fmt.Errorf("Error fetching tweets: %v", err)
	}
	defer rows.Close()

	tweets := []models.Tweet{}

	for rows.Next() {
		var tweet models.Tweet

		err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Content, &tweet.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}

		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating rows: %v", err)
	}

	return tweets, nil
}

// Funcion para obtener el timeline de los usuarios a los que se sigue
func GetTweetsTimeline(db *sql.DB, userId *int64) ([]models.Tweet, error) {
	query := `SELECT tw.id as tw_id, tw.user_id, us.name, tw.content, tw.created_at as tweet_date
				FROM tweets AS tw
				INNER JOIN follows AS fol ON fol.followed_id = tw.user_id
				INNER JOIN users AS us ON us.id = tw.user_id
				WHERE fol.follower_id = $1
				ORDER BY tweet_date DESC;`
	rows, err := db.Query(query, userId)

	if err != nil {
		return nil, fmt.Errorf("Error fetching timeline: %v", err)
	}

	timeline := []models.Tweet{}

	for rows.Next() {
		var tweet models.Tweet

		err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.AuthorName, &tweet.Content, &tweet.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}

		timeline = append(timeline, tweet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating rows: %v", err)
	}

	return timeline, nil
}
