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

type TweetController struct {
	TweetService *services.TweetService
}

func NewTweetController(db *sql.DB) *TweetController {
	tweetService := services.NewTweetService(db)
	return &TweetController{TweetService: tweetService}
}

func (tc *TweetController) CreateTweetHandler(c *gin.Context) {
	var tweet models.Tweet

	if err := c.ShouldBindJSON(&tweet); err != nil {
		errorResponse := utils.ResponseToApi(http.StatusInternalServerError, err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	tweetPosted, err := tc.TweetService.PostTweet(&tweet)
	if err != nil {
		errorResponse := utils.ResponseToApi(http.StatusInternalServerError, "[X] Error posting tweet: "+err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if !tweetPosted {
		errorResponse := utils.ResponseToApi(http.StatusInternalServerError, "[X] Could not post the tweet", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := utils.ResponseToApi(http.StatusCreated, "Tweet posted", false, 0, 0, 0)

	c.JSON(http.StatusCreated, response)
}

func (tc *TweetController) GetTimelineHandler(c *gin.Context) {
	idStr := c.Param("follower_id")

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "Invalid follower ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	// Llamamos al servicio para obtener el usuario
	timeline, err := tc.TweetService.GetUserTimeline(&id)

	if err != nil {
		errorResponse := utils.ResponseToApi(http.StatusInternalServerError, err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	response := utils.ResponseToApi(http.StatusOK, timeline, false, 0, 0, 0)
	c.JSON(http.StatusOK, response)
}
