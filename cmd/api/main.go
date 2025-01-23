package main // Nombre del paquete

import (
	"net/http"

	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/db"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/utils"
	"github.com/gin-gonic/gin"
) // Importar dependencias

func main() { // función inicio requerida
	db.SetupSQLiteDatabase()

	router := gin.Default()

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

	router.Run(":8080")
}
