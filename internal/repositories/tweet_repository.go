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
func GetTweetsTimeline(db *sql.DB, redisClient *redis.Client, userId *int64, limit *int64, offset *int64) ([]models.Tweet, error) {
	var ctx = context.Background()
	cacheKey := fmt.Sprintf("timeline:%d:%d:%d", *userId, *limit, *offset)
	if redisClient != nil {
		// Crear una clave de cache única basada en userId, limit y offset

		fmt.Println("[x] Buscando data en Redis en la key %v", cacheKey)
		var cachedTimeline []models.Tweet

		cachedTimelineData, err := redisClient.Get(ctx, cacheKey).Result()

		//Si por alguna razon falla la obtención de datos a través de Redis, la ejecucion no se detiene y se intenta obtener la data a traves de la db sql
		if err == nil {
			// Si hay datos en cache, deserializarlos y devolver
			err = json.Unmarshal([]byte(cachedTimelineData), &cachedTimeline)
			if err != nil {
				//Ideal: Generar log de auditoria para saber porque falló la deserialización
				fmt.Println("Error deserializing Redis data: %v", err)
			} else {
				fmt.Println("Data founded in Redis!")
				// Si se encuentra en cache, devolver el resultado
				return cachedTimeline, nil
			}

		} else if err != redis.Nil {
			// Si hay otro tipo de error, lo manejamos
			//Ideal: Generar log de auditoria para saber porque falló redis
			fmt.Println("Error to get Redis data: %v", err)
		}
	}

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
		return nil, fmt.Errorf("Error fetching timeline: %v", err)
	}
	defer rows.Close()

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

	if redisClient != nil {
		//Preparo el time line en formato json para guardarlos en Redis
		timelineJSON, err := json.Marshal(timeline)
		if err != nil {
			fmt.Println("Error serializing data for Redis: %v", err)
		}

		//Establezco ttl de 1 hora, pasado ese tiempo, la data se borra del cache
		err = redisClient.Set(ctx, cacheKey, timelineJSON, 1*time.Hour).Err()
		if err != nil {
			fmt.Println("Error to set value in Redis: %v", err)
		}
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
