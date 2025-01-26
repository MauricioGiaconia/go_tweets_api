package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/repositories"
	"github.com/redis/go-redis/v9"
)

type TweetService struct {
	DB  *sql.DB       // Conexion a db SQL
	RDB *redis.Client // Conexion a db redis
}

func NewTweetService(db *sql.DB, rdb *redis.Client) *TweetService {
	return &TweetService{DB: db, RDB: rdb}
}

func (ts *TweetService) GetUserTimeline(followerId *int64, limit *int64, offset *int64) ([]models.Tweet, error) {

	cacheKey := fmt.Sprintf("timeline:%d:%d:%d", *followerId, *limit, *offset)

	if ts.RDB != nil {
		cachedTimeline, err := repositories.GetTweetsFromCache(ts.RDB, cacheKey)
		if err != nil {
			fmt.Println("Error getting timeline from Redis: %v", err) // No detengo la ejecución asi se intenta obtener la data solicitada desde la DB sql
		}

		// Si los datos están en cache, los devolvemos
		if cachedTimeline != nil {
			if cachedTimeline.IsFullPage {
				fmt.Println("[x] Returning data from cache!")
				return cachedTimeline.Tweets, nil
			}
			fmt.Println("[x] The timeline consulted may be outdated, searching for information in the sql database...")
		}
	}

	timeline, err := repositories.GetTweetsFromDB(ts.DB, followerId, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("Error getting timeline: %v", err)
	}

	//Se guarda unicamente en redis si hay informacion
	if int64(len(timeline)) > 0 {
		if ts.RDB != nil {
			isFullPage := int64(len(timeline)) == *limit

			timelineCache := models.TimelineCache{
				Tweets:     timeline,
				IsFullPage: isFullPage,
			}

			ttl := 30 * time.Minute

			if !isFullPage {
				//En caso que la pagina NO este completa, se mantiene un time to live menor
				ttl = 10 * time.Minute
			}

			err = repositories.SaveTweetsToCache(ts.RDB, cacheKey, &timelineCache, ttl)
			if err != nil {
				fmt.Println("Error saving timeline to Redis: %v", err) // Si no se pudo guardar la data en cache, retorno de todas formas la informacion obtenida de la db sql
			}
		}
	} else {
		timeline = []models.Tweet{}
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
	//El len se hace sobre rune para tratar de forma correct a los caracteres multibtyes (como acentos, simbolos etc etc)
	if len([]rune(tweet.Content)) > maxCharacters {
		return false, fmt.Errorf("The content of the tweet must not exceed 280 characters")
	}

	_, err := repositories.GetUserById(ts.DB, tweet.UserID)

	if err != nil {
		return false, fmt.Errorf("Nonexistent user")
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
