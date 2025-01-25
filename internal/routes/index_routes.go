// routes/index_routes.go
package routes

import (
	"database/sql"
	"net/http"

	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Todas las rutas estaran centralizadas en SetupRoutes
func SetupRoutes(router *gin.Engine, db *sql.DB) {

	// Rutas relacionadas con usuarios
	SetupUserRoutes(router, db)

	// Rutas relacionadas con seguidores
	SetupUserFollowRoutes(router, db)

	// Rutas relacionadas con tweets
	SetupTweetRoutes(router, db)

	//Endpoint ping para probar el funcionamiento de la API
	router.GET("/ping", func(c *gin.Context) {
		response := utils.ResponseToApi(http.StatusOK, "Pong", false, 0, 0, 0)

		c.JSON(http.StatusOK, response)
	})

	// Maneja rutas no encontradas
	router.NoRoute(func(c *gin.Context) {
		notFoundResponse := utils.ResponseToApi(http.StatusNotFound, "Endpoint not found", false, 0, 0, 0)
		c.JSON(http.StatusNotFound, notFoundResponse)
	})

}
