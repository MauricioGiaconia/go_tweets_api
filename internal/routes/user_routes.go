package routes

import (
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configura las rutas para manejar usuarios.
func SetupUserRoutes(router *gin.Engine) {

	userGroup := router.Group("/users")
	{
		userGroup.POST("/", controllers.CreateUserHandler) // POST /users crea un nuevo usuario
	}
}
