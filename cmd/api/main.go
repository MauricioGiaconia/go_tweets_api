package main

import (
	"flag"
	"log"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/routes"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/db"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/factory"
	"github.com/gin-gonic/gin"
)

func main() {

	//Obtengo el tipo de db a utilizar por linea de comandos, usando la flag -db
	dbType := flag.String("db", "sqlite", "Tipo de base de datos a usar (postgres, sqlite)")
	flag.Parse()

	if *dbType != "sqlite" && *dbType != "postgres" {
		log.Fatalf("[x] Invalid dbType")
	}

	dbInstance, err := factory.GetDatabase(*dbType)

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
