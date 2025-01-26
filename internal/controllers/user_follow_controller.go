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
	"github.com/redis/go-redis/v9"
)

type UserFollowController struct {
	UserFollowService *services.FollowService
}

func NewUseFollowrController(db *sql.DB, rds *redis.Client) *UserFollowController {
	userFollowService := services.NewFollowService(db, rds)
	return &UserFollowController{UserFollowService: userFollowService}
}

// FollowUserHandler maneja la solicitud de seguimiento de un usuario a otro
func (ufc *UserFollowController) FollowUserHandler(c *gin.Context) {
	var follow models.UserFollow

	// Decodificamos el cuerpo de la solicitud JSON al struct User
	if err := c.ShouldBindJSON(&follow); err != nil {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "Error decoding body", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	if follow.FollowerID == follow.FollowedID {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "Cannot follow yourself", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	if follow.FollowedID <= 0 || follow.FollowedID <= 0 {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "Invalid follower or followed ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	// Llamamos al servicio para crear el usuario
	followResponse, err := ufc.UserFollowService.FollowUser(&follow)

	if err != nil {
		if err.Error() == "Nonexistent followed ID user" || err.Error() == "Nonexistent follower ID user" {
			badResponse := utils.ResponseToApi(http.StatusBadRequest, err.Error(), false, 0, 0, 0)
			c.JSON(http.StatusBadRequest, badResponse)
			return
		}

		serverErrorResponse := utils.ResponseToApi(http.StatusInternalServerError, err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, serverErrorResponse)
		return
	}

	msgResponse := "Followed"
	responseCode := http.StatusCreated
	if !followResponse {
		responseCode = http.StatusBadRequest
		msgResponse = "Cannot follow the user"
	}

	finalResponse := utils.ResponseToApi(int64(responseCode), msgResponse, false, 0, 0, 0)

	c.JSON(responseCode, finalResponse)
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
