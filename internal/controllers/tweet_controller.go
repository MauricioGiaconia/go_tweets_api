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

type TweetController struct {
	TweetService *services.TweetService
}

func NewTweetController(db *sql.DB, rdb *redis.Client) *TweetController {
	tweetService := services.NewTweetService(db, rdb)
	return &TweetController{TweetService: tweetService}
}

func (tc *TweetController) CreateTweetHandler(c *gin.Context) {
	var tweet models.Tweet

	if err := c.ShouldBindJSON(&tweet); err != nil {
		errorResponse := utils.ResponseToApi(http.StatusInternalServerError, err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if tweet.UserID <= 0 {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "Invalid user ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, badResponse)
		return
	}

	tweetPosted, err := tc.TweetService.PostTweet(&tweet)

	if err != nil {

		if err.Error() == "Nonexistent user" || err.Error() == "The content of the tweet must not exceed 280 characters" {
			badResponse := utils.ResponseToApi(http.StatusBadRequest, err.Error(), false, 0, 0, 0)
			c.JSON(http.StatusBadRequest, badResponse)
			return
		}

		errorResponse := utils.ResponseToApi(http.StatusInternalServerError, "[X] Error posting tweet: "+err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if !tweetPosted {
		errorResponse := utils.ResponseToApi(http.StatusInternalServerError, "[X] Could not post the tweet", false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	response := utils.ResponseToApi(http.StatusCreated, "Tweet posted", false, 0, 0, 0)

	c.JSON(http.StatusCreated, response)
}

func (tc *TweetController) GetTimelineHandler(c *gin.Context) {
	idStr := c.Param("follower_id")

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil || id <= 0 {
		badResponse := utils.ResponseToApi(http.StatusBadRequest, "Invalid follower ID", false, 0, 0, 0)
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

	timeline, err := tc.TweetService.GetUserTimeline(&id, &limit, &offset)

	if err != nil {
		if err.Error() == "Nonexistent user" {
			badResponse := utils.ResponseToApi(http.StatusBadRequest, err.Error(), false, 0, 0, 0)
			c.JSON(http.StatusBadRequest, badResponse)
			return
		}

		errorResponse := utils.ResponseToApi(http.StatusInternalServerError, err.Error(), false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	totalTweets, err := tc.TweetService.CountTimeline(&id)

	if err != nil {
		//Ideal: Implementar creacion de log indicando cual fue el error en el count
		fmt.Println(err)
	}

	//Por mas que el count rompa, retorno la informacion igual ya que cuento con el timeline
	response := utils.ResponseToApi(http.StatusOK, timeline, true, totalTweets, limit, offset)
	c.JSON(http.StatusOK, response)
}
