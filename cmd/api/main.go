package main

import (
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// Configurar las rutas
	routes.SetupRoutes(router) // Llamar a las rutas de usuarios

	router.Run(":8080")
}
