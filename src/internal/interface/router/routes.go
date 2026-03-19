package router

import (
	"wetalk-academy/config"
	repository "wetalk-academy/internal/infrastructure/db/repository"
	"wetalk-academy/internal/interface/handler"
	"wetalk-academy/internal/interface/middleware"
	"wetalk-academy/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppHandler struct {
	topicHandler  *handler.TopicHandler
	lessonHandler *handler.LessonHandler
}

func SetupRoutes(mongoDB *mongo.Database, conf *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware(conf.App.Whitelist))

	// repositories
	topicRepo := repository.NewTopicRepository(mongoDB)
	lessonRepo := repository.NewLessonRepository(mongoDB)

	// services
	topicService := service.NewTopicService(topicRepo, lessonRepo)
	lessonService := service.NewLessonService(lessonRepo, topicRepo)

	// handlers
	appHandler := &AppHandler{
		topicHandler:  handler.NewTopicHandler(topicService),
		lessonHandler: handler.NewLessonHandler(lessonService),
	}

	api := router.Group("/api/v1")
	{
		setupPublicRoutes(api, appHandler)
		setupProtectedRoutes(api, conf, appHandler)
	}

	return router
}

func setupPublicRoutes(rg *gin.RouterGroup, appHandler *AppHandler) {
	rg.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	topics := rg.Group("/topics")
	{
		topics.GET("", appHandler.topicHandler.GetTopics)
		topics.GET("/:slug", appHandler.topicHandler.GetTopicBySlug)
		topics.GET("/:slug/lessons", appHandler.topicHandler.GetLessonsInTopic)
	}

	lessons := rg.Group("/lessons")
	{
		lessons.GET("/:slug", appHandler.lessonHandler.GetLessonBySlug)
	}
}

func setupProtectedRoutes(rg *gin.RouterGroup, conf *config.Config, appHandler *AppHandler) {
	protected := rg.Group("")
	protected.Use(middleware.AuthMiddleware(conf))
	{
		topics := protected.Group("/topics")
		{
			topics.POST("", appHandler.topicHandler.CreateTopic)
			topics.PUT("/:slug", appHandler.topicHandler.UpdateTopic)
			topics.DELETE("/:slug", appHandler.topicHandler.DeleteTopic)
		}

		lessons := protected.Group("/lessons")
		{
			lessons.POST("", appHandler.lessonHandler.CreateLesson)
			lessons.PUT("/:slug", appHandler.lessonHandler.UpdateLesson)
			lessons.DELETE("/:slug", appHandler.lessonHandler.DeleteLesson)
		}
	}
}
