package factory

import (
	"fmt"

	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/db"
	redisdb "github.com/MauricioGiaconia/uala_backend_challenge/pkg/redis_db"
	"github.com/redis/go-redis/v9"
)

// GetDatabase crea una instancia de la base de datos que se solicite (SQLite o PostgreSQL)
func GetDatabase(dbType string) (db.Database, error) {
	switch dbType {
	case "sqlite":
		return &db.SQLiteDatabase{}, nil
	case "postgres":
		return &db.PostgresDatabase{}, nil
	default:
		return nil, fmt.Errorf("[x] Invalid database type: %s", dbType)
	}
}

// GetCache crea una instancia del cache solicitado (Redis)
func GetCache(cacheType string) (*redis.Client, error) {
	switch cacheType {
	case "redis":
		client, err := redisdb.NewRedisClient("tes", "test", 0)
		if err != nil {
			return nil, err
		}
		return client, nil
	default:
		return nil, fmt.Errorf("[x] Invalid cache type: %s", cacheType)
	}
}
