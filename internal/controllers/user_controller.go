package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(db *sql.DB) *UserController {
	// Se inicializa el servicio de usuarios pasandole la instancia de la DB
	userService := services.NewUserService(db)
	return &UserController{UserService: *userService}
}

func (uc *UserController) CreateUserHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("[x] Error decoding body: %v", err),
		})
		return
	}

	userID, err := uc.UserService.CreateUser(&user)
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

func (uc *UserController) GetUserByIdHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := uc.UserService.GetUserById(id)
	fmt.Println("------------------------")
	fmt.Println(err)
	fmt.Println("------------------------")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Can't get user",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data": user,
	})
}
