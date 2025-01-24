package main

import (
	"log"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/routes"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/db"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/factory"
	"github.com/gin-gonic/gin"
)

func main() {

	dbInstance, err := factory.GetDatabase("postgres") //Esta hardcodeado sqlite, lo ideal seria recibir por variable de entorno la base de datos a utilziar

	if err != nil {
		log.Fatalf("[x] Error getting database instance: %v", err)
	}

	conn, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("[x] Error connecting to database: %v", err)
	}

	router := gin.Default()

	// Configurar las rutas
	routes.SetupRoutes(router, conn)

	defer db.CloseDatabase(conn)

	router.Run(":8080")
}
