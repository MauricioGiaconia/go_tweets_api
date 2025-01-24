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

type UserFollowController struct {
	UserFollowService *services.FollowService
}

func NewUseFollowrController(db *sql.DB) *UserFollowController {
	userFollowService := services.NewFollowService(db)
	return &UserFollowController{UserFollowService: userFollowService}
}

// FollowUserHandler maneja la solicitud de seguimiento de un usuario a otro
func (ufc *UserFollowController) FollowUserHandler(c *gin.Context) {
	var follow models.UserFollow

	// Decodificamos el cuerpo de la solicitud JSON al struct User
	if err := c.ShouldBindJSON(&follow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("[x] Error decoding body: %v", err),
		})
		return
	}

	if follow.FollowerID == follow.FollowedID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot follow yourself",
		})
		return
	}

	// Llamamos al servicio para crear el usuario
	followResponse, err := ufc.UserFollowService.FollowUser(&follow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("[x] Error to follow: %v", err),
		})
		return
	}

	msgResponse := "Followed"

	if !followResponse {
		msgResponse = "Cannot follow the user"
	}

	// Respondemos con el ID del usuario creado
	c.JSON(http.StatusCreated, gin.H{
		"msg": msgResponse,
	})
}

// GetFollowersHandler maneja la solicitud de obtener un usuario por su ID.
func (ufc *UserFollowController) GetFollowersHandler(c *gin.Context) {
	idStr := c.Param("id")
	relationType := c.Param("follow_type")

	// Validar que el parámetro follow_type sea válido
	if relationType != "followers" && relationType != "following" {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "Invalid follow type. Must be 'followers' or 'following'", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "Invalid user ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	// Llamamos al servicio para obtener el usuario
	userFollowInfo, err := ufc.UserFollowService.GetFollows(&id, &relationType)

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

	// Respondemos con los seguidores/seguidos del usuario en formato JSON junto a la informacion del paginado
	response := utils.ResponseToApi(http.StatusOK, userFollowInfo, true, 0, 0, 0)
	c.JSON(http.StatusOK, response)
}
