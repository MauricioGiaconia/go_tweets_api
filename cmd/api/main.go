package main

import (
	"flag"
	"fmt"
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

	dbConn, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("[x] Error connecting to database: %v", err)
	}

	redisClient, err := factory.GetCache("redis")
	if err != nil {
		fmt.Println("API working without redis: %v", err.Error())
	} else {
		fmt.Println("API working with redis!")
		defer redisClient.Close()
	}

	router := gin.Default()

	// Configurar las rutas
	routes.SetupRoutes(router, dbConn, redisClient)

	defer db.CloseDatabase(dbConn)

	router.Run(":8080")
}
