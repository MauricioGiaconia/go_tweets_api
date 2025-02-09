package routes

import (
	"database/sql"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupTweetRoutes(router *gin.Engine, db *sql.DB, rdb *redis.Client) {

	tweetController := controllers.NewTweetController(db, rdb)

	tweetGroup := router.Group("/tweets")
	{
		tweetGroup.POST("/create", tweetController.CreateTweetHandler)                                    // POST /tweets/post crea un nuevo tweet
		tweetGroup.GET("/:follower_id/timeline", tweetController.GetTimelineHandler)                      // GET /tweets/:follower_id/timeline obtengo el timeline de los usuarios seguidos
		tweetGroup.GET("/:follower_id/routine_timeline", tweetController.GetTimelineWithGoRoutineHandler) // GET /tweets/:follower_id/routine_timeline obtengo el timeline de los usuarios seguidos usango go routines
	}
}
