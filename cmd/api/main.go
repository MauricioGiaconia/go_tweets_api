package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/routes"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/db"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/factory"
	"github.com/gin-gonic/gin"
)

func main() {

	//Obtengo el tipo de db y port a utilizar por linea de comandos, usando la flag -db y -port
	dbType := flag.String("db", "sqlite", "Tipo de base de datos a usar (postgres, sqlite)")
	port := flag.String("port", "8080", "Puerto a utilizar") //Por defecto se usa el puerto 8080
	flag.Parse()

	portNum, err := strconv.Atoi(*port)
	if err != nil || portNum < 1 || portNum > 65535 {
		log.Fatalf("[x] Invalid port")
	}

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

	address := fmt.Sprintf(":%s", *port)
	if err := router.Run(address); err != nil {
		log.Fatalf("[x] Failed to start server: %v", err)
	}
}
