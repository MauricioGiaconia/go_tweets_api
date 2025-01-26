package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/redis/go-redis/v9"
)

// Funciones para interactura con db SQL
func PostTweet(db *sql.DB, tweet *models.Tweet) (bool, error) {
	tx, err := db.Begin() //Se inicia transaccion para ejecutar Rollback si algo sale mal
	query := `INSERT INTO tweets (user_id, content) VALUES ($1, $2)`

	_, err = tx.Exec(query, tweet.UserID, tweet.Content)
	if err != nil {
		tx.Rollback()
		fmt.Println("[x] Error to create Tweet: %v", err)
		return false, fmt.Errorf(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return false, fmt.Errorf("Error committing transaction: %v", err)
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
func GetTweetsFromDB(db *sql.DB, userId *int64, limit *int64, offset *int64) ([]models.Tweet, error) {
	query := `SELECT tw.id as tw_id, tw.user_id, us.name, tw.content, tw.created_at as tweet_date
              FROM tweets AS tw
              INNER JOIN follows AS fol ON fol.followed_id = tw.user_id
              INNER JOIN users AS us ON us.id = tw.user_id
              WHERE fol.follower_id = $1
              ORDER BY tweet_date DESC
              LIMIT $2
              OFFSET $3;`
	rows, err := db.Query(query, userId, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("Error fetching timeline from DB: %v", err)
	}
	defer rows.Close()

	var timeline []models.Tweet
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

func CountTweetsTimeline(db *sql.DB, userId *int64) (int64, error) {
	query := `SELECT COUNT(*) AS total_tweets
				FROM tweets AS tw
				INNER JOIN follows AS fol ON fol.followed_id = tw.user_id
				INNER JOIN users AS us ON us.id = tw.user_id
				WHERE fol.follower_id = $1;`

	rows, err := db.Query(query, userId)
	if err != nil {
		return 0, fmt.Errorf("Error fetching timeline: %v", err)
	}
	defer rows.Close()

	var totalTweets int64

	if rows.Next() {
		err = rows.Scan(&totalTweets)
		if err != nil {
			return 0, fmt.Errorf("Error scanning row: %v", err)
		}
	} else {

		return 0, fmt.Errorf("No rows found")
	}

	return totalTweets, nil
}

//Funciones para interactuar con redis respecto a los Tweets

func GetTweetsFromCache(redisClient *redis.Client, cacheKey string) (*models.TimelineCache, error) {
	var ctx = context.Background()
	cachedTimelineData, err := redisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil // No hay datos en cache
	} else if err != nil {
		return nil, fmt.Errorf("Error getting data from Redis: %v", err)
	}

	var cachedTimeline models.TimelineCache
	err = json.Unmarshal([]byte(cachedTimelineData), &cachedTimeline)
	if err != nil {
		return nil, fmt.Errorf("Error deserializing Redis data: %v", err)
	}

	return &cachedTimeline, nil
}

func SaveTweetsToCache(redisClient *redis.Client, cacheKey string, timeline *models.TimelineCache, ttl time.Duration) error {
	var ctx = context.Background()
	timelineJSON, err := json.Marshal(timeline)
	if err != nil {
		return fmt.Errorf("Error serializing data for Redis: %v", err)
	}

	err = redisClient.Set(ctx, cacheKey, timelineJSON, ttl).Err()
	if err != nil {
		return fmt.Errorf("Error setting value in Redis: %v", err)
	}

	return nil
}
