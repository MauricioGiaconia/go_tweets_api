package factory

import (
	"fmt"

	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/db"
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

// // GetCache crea una instancia del cache solicitado (Redis)
// func GetCache(cacheType string) (interface{}, error) {
// 	switch cacheType {
// 	case "redis":
// 		return &cache.RedisCache{}, nil
// 	default:
// 		return nil, fmt.Errorf("[x] Invalid cache type: %s", cacheType)
// 	}
// }
