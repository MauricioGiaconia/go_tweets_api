package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/models"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/services"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(db *sql.DB) *UserController {
	userService := services.NewUserService(db)
	return &UserController{UserService: userService}
}

// CreateUserHandler maneja la solicitud de creación de un nuevo usuario.
func (uc *UserController) CreateUserHandler(c *gin.Context) {
	var user models.User

	// Decodificamos el cuerpo de la solicitud JSON al struct User
	if err := c.ShouldBindJSON(&user); err != nil {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "[x] Error decoding body: "+err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	// Llamamos al servicio para crear el usuario
	userID, err := uc.UserService.CreateUser(&user)
	if err != nil {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "[x] Error creating user: "+err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	// Respondemos con el ID del usuario creado
	response := utils.ResponseToApi(http.StatusCreated, gin.H{
		"id": userID,
	}, false, 0, 0, 0)
	c.JSON(http.StatusCreated, response)
}

// GetUserByIdHandler maneja la solicitud de obtener un usuario por su ID.
func (uc *UserController) GetUserByIdHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil || id <= 0 {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "Invalid user ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	// Llamamos al servicio para obtener el usuario
	user, err := uc.UserService.GetUserById(id)

	if err != nil {
		if err.Error() == "Error fetching user: user not found" {
			notFoundResponse := utils.ResponseToApi(404, "Not found", false, 0, 0, 0)
			c.JSON(404, notFoundResponse)
			return
		}

		errorResponse := utils.ResponseToApi(http.StatusInternalServerError, err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	// Respondemos con los datos del usuario en formato JSON
	response := utils.ResponseToApi(http.StatusOK, user, false, 0, 0, 0)
	c.JSON(http.StatusOK, response)
}
