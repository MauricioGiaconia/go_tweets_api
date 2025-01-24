package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/services"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/utils"
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
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "Invalid user ID", false, 0, 0, 0)

		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	user, err := uc.UserService.GetUserById(id)

	if err != nil {

		if err.Error() == "user not found" {
			notFoundResponse := utils.ResponseToApi(404, err.Error(), false, 0, 0, 0)
			c.JSON(404, notFoundResponse)
			return
		}

		errorResponse := utils.ResponseToApi(http.StatusInternalServerError, err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	response := utils.ResponseToApi(http.StatusOK, user, false, 0, 0, 0)

	c.JSON(http.StatusOK, response)
}
