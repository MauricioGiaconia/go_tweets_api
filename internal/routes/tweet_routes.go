package routes

import (
	"database/sql"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupTweetRoutes(router *gin.Engine, db *sql.DB) {

	tweetController := controllers.NewTweetController(db)

	tweetGroup := router.Group("/tweets")
	{
		tweetGroup.POST("/create", tweetController.CreateTweetHandler)               // POST /tweets/post crea un nuevo tweet
		tweetGroup.GET("/:follower_id/timeline", tweetController.GetTimelineHandler) // GET /tweets/:follower_id/timeline obtengo el timeline de los usuarios seguidos
	}
}
