package test

import (
	"database/sql"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func SetupRouter(db *sql.DB, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	routes.SetupRoutes(router, db, rdb)
	return router
}
