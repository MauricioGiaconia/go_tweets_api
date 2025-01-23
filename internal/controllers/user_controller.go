package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/db"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/services"
	"github.com/gin-gonic/gin"
)

func CreateUserHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("[x] Error decoding body: %v", err),
		})
		return
	}

	database, err := db.SetupSQLiteDatabase() // En caso de utiliza postgresql, reemplazar con el metodo que corresponde
	if err != nil {
		log.Println("[x] CreateUserHandler DB connection error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "CreateUserHandler DB connection error",
		})
		return
	}
	// En caso de utilziar postgress, descomentar el cierre de conexion a la DB.
	// Si se cierra la conexion de SQLite utilizado en memoria, los datos se pierden.
	// defer db.CloseDatabase(database)

	userID, err := services.CreateUserService(database, &user)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("[x] Error creating user: %v", err),
		})
		return
	}

	// Respondo con el ID del usuario creado
	c.JSON(http.StatusCreated, gin.H{
		"id": userID,
	})
}
