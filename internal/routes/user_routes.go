package routes

import (
	"database/sql"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configura las rutas para manejar usuarios.
func SetupUserRoutes(router *gin.Engine, db *sql.DB) {

	userController := controllers.NewUserController(db)

	userGroup := router.Group("/users")
	{
		userGroup.POST("/create", userController.CreateUserHandler) // POST /users crea un nuevo usuario
	}
}
