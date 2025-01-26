package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/repositories"
	"github.com/redis/go-redis/v9"
)

type FollowService struct {
	DB  *sql.DB
	RDB *redis.Client
}

func NewFollowService(db *sql.DB, rdb *redis.Client) *FollowService {
	return &FollowService{DB: db, RDB: rdb}
}

func (ufs *FollowService) FollowUser(follow *models.UserFollow) (bool, error) {

	_, err := repositories.GetUserById(ufs.DB, follow.FollowedID)

	if err != nil {
		return false, fmt.Errorf("Nonexistent followed ID user")
	}

	_, err = repositories.GetUserById(ufs.DB, follow.FollowerID)

	if err != nil {
		return false, fmt.Errorf("Nonexistent follower ID user")
	}

	userFollow, err := repositories.FollowUser(ufs.DB, follow)

	if err != nil {
		return userFollow, fmt.Errorf("Error followed user: %v", err)
	}

	return userFollow, nil
}

func (ufs *FollowService) GetFollows(userId *int64, relationType *string, limit *int64, offset *int64) (models.UserFollows, error) {

	// Validar que el relationType sea el adecuado segun la logica implementada en el repository
	if *relationType != "followers" && *relationType != "following" {

		return models.UserFollows{}, fmt.Errorf("Invalid follow type. Must be 'followers' or 'following'")
	}

	cacheKey := fmt.Sprintf("follows:%d:%d:%d:%d", *userId, *relationType, *limit, *offset)

	if ufs.RDB != nil {
		cachedFollows, err := repositories.GetFollowsFromCache(ufs.RDB, cacheKey)
		if err != nil {
			fmt.Println("Error getting follows from Redis: %v", err) // No detengo la ejecución asi se intenta obtener la data solicitada desde la DB sql
		}

		// Si los datos están en cache, los devolvemos
		if cachedFollows != nil {
			if cachedFollows.IsFullPage {
				fmt.Println("[x] Returning data from cache!")
				return cachedFollows.Follows, nil
			}
			fmt.Println("[x] The timeline consulted may be outdated, searching for information in the sql database...")
		}
	}

	userFollows, err := repositories.GetFollows(ufs.DB, *userId, *relationType, limit, offset)

	if err != nil {
		return models.UserFollows{}, fmt.Errorf("Error getting follows: %v", err)
	}

	if ufs.RDB != nil {
		isFullPage := int64(len(userFollows.Follows)) == *limit

		followsCache := models.FollowsCache{
			Follows:    *userFollows,
			IsFullPage: isFullPage,
		}

		ttl := 30 * time.Minute

		if !isFullPage {
			//En caso que la pagina NO este completa, se mantiene un time to live menor
			ttl = 10 * time.Minute
		}

		err = repositories.SaveFollowsToCache(ufs.RDB, cacheKey, &followsCache, ttl)
		if err != nil {
			fmt.Println("Error saving timeline to Redis: %v", err) // Si no se pudo guardar la data en cache, retorno de todas formas la informacion obtenida de la db sql
		}
	}
	return *userFollows, nil
}

func (ufs *FollowService) CountFollows(userId *int64, relationType *string) (int64, error) {
	total, err := repositories.CountFollows(ufs.DB, *userId, *relationType)

	if err != nil {
		return 0, fmt.Errorf("Error counting timeline: %v", err)
	}

	return total, nil
}
