package routes

import (
	"database/sql"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// SetupUserRoutes configura las rutas para manejar usuarios.
func SetupUserFollowRoutes(router *gin.Engine, db *sql.DB, rds *redis.Client) {

	userFollowController := controllers.NewUseFollowrController(db, rds)

	userFollowGroup := router.Group("/users_follow")
	{
		userFollowGroup.POST("/create", userFollowController.FollowUserHandler)                    // POST /users_follow/create crea un nuevo usuario seguidor a un usuario
		userFollowGroup.GET("/:id/follows/:follow_type", userFollowController.GetFollowersHandler) // GET /users_follow/:id/followers obtiene todos los seguidores de un usuario
	}
}
