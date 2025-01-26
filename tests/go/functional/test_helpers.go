package functional

import (
	"database/sql"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func setupRouter(db *sql.DB, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	routes.SetupRoutes(router, db, rdb)
	return router
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type SuccessResponse struct {
	Code int   `json:"code"`
	Data int64 `json:"data"`
}
