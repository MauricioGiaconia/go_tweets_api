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

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "Invalid user ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	const defaultLimit int64 = 25 // Por defecto, vendran 25 tweets por pagina
	const defaultOffset int64 = 0
	const maxLimit int64 = 100 // Limite máximo permitido

	var limit, offset int64
	var paramError error

	if limitStr != "" {
		limit, paramError = strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limit <= 0 || limit > maxLimit {
			badResponse := utils.ResponseToApi(http.StatusBadRequest, "Invalid limit parameter", false, 0, 0, 0)
			c.JSON(http.StatusBadRequest, badResponse)
			return
		}
	} else {
		limit = defaultLimit
	}

	if offsetStr != "" {
		offset, paramError = strconv.ParseInt(offsetStr, 10, 64)
		if paramError != nil || offset < 0 {
			badResponse := utils.ResponseToApi(http.StatusBadRequest, "Invalid offset parameter", false, 0, 0, 0)
			c.JSON(http.StatusBadRequest, badResponse)
			return
		}
	} else {
		offset = defaultOffset
	}

	// Llamamos al servicio para obtener el usuario
	userFollowInfo, err := ufc.UserFollowService.GetFollows(&id, &relationType, &limit, &offset)

	if err != nil {
		errorResponse := utils.ResponseToApi(http.StatusInternalServerError, err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	totalFollows, err := ufc.UserFollowService.CountFollows(&id, &relationType)

	if err != nil {
		//Ideal: Implementar creacion de log indicando cual fue el error en el count
		fmt.Println(err)
	}

	// Respondemos con los seguidores/seguidos del usuario en formato JSON junto a la informacion del paginado
	response := utils.ResponseToApi(http.StatusOK, userFollowInfo, true, totalFollows, limit, offset)
	c.JSON(http.StatusOK, response)
}
