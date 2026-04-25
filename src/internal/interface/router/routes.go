package router

import (
	"strings"
	"wetalk-academy/config"
	"wetalk-academy/internal/interface/handler"
	"wetalk-academy/internal/interface/middleware"
	"wetalk-academy/internal/wire"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(appHandler *wire.AppHandler, conf *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware(conf.App.Whitelist))

	api := router.Group("/api/v1")
	{
		api.Use(middleware.RequestMetricsMiddleware())
		setupPublicRoutes(api, conf, appHandler)
		setupProtectedRoutes(api, conf, appHandler)
	}

	if strings.TrimSpace(conf.Log.DashboardToken) != "" {
		dash := router.Group("")
		dash.Use(middleware.LogDashboardAuth(conf))
		dash.GET("/admin/logs", handler.ServeAdminLogsDashboard)

		apiAdmin := router.Group("/api/v1/admin")
		apiAdmin.Use(middleware.LogDashboardAuth(conf))
		apiAdmin.GET("/logs/files", handler.GetAdminLogFiles)
		apiAdmin.GET("/logs", handler.GetAdminLogs)
		apiAdmin.GET("/metrics", handler.GetAdminMetrics)
	}

	return router
}

func setupPublicRoutes(rg *gin.RouterGroup, conf *config.Config, appHandler *wire.AppHandler) {
	rg.Use(middleware.OptionalAuthMiddleware(conf))

	rg.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	topics := rg.Group("/topics")
	{
		topics.GET("", appHandler.TopicHandler.GetTopics)
		topics.GET("/:slug", appHandler.TopicHandler.GetTopicBySlug)
		topics.GET("/:slug/lessons", appHandler.TopicHandler.GetLessonsInTopic)
	}

	lessons := rg.Group("/lessons")
	{
		lessons.GET("/:slug", appHandler.LessonHandler.GetLessonBySlug)
		lessons.GET("/:slug/content", appHandler.ContentHandler.GetContent)
		lessons.GET("/:slug/quiz", appHandler.QuizHandler.GetQuizzesByLessonSlug)
	}

	quizzes := rg.Group("/quizzes")
	{
		quizzes.GET("/:id", appHandler.QuizHandler.GetQuizByID)
	}

	judge0 := rg.Group("/judge0")
	{
		judge0.POST("/submit", appHandler.Judge0Handler.SubmitCode)
	}
}

func setupProtectedRoutes(rg *gin.RouterGroup, conf *config.Config, appHandler *wire.AppHandler) {
	protected := rg.Group("")
	protected.Use(middleware.AuthMiddleware(conf))
	{
		topics := protected.Group("/topics")
		{
			topics.POST("", appHandler.TopicHandler.CreateTopic)
			topics.PUT("/:slug", appHandler.TopicHandler.UpdateTopic)
			topics.DELETE("/:slug", appHandler.TopicHandler.DeleteTopic)
		}

		lessons := protected.Group("/lessons")
		{
			lessons.POST("", appHandler.LessonHandler.CreateLesson)
			lessons.PUT("/:slug", appHandler.LessonHandler.UpdateLesson)
			lessons.DELETE("/:slug", appHandler.LessonHandler.DeleteLesson)

			lessons.POST("/:slug/content", appHandler.ContentHandler.CreateContent)
			lessons.PUT("/:slug/content", appHandler.ContentHandler.UpdateContent)
			lessons.DELETE("/:slug/content", appHandler.ContentHandler.DeleteContent)
		}

		quizzes := protected.Group("/quizzes")
		{
			quizzes.POST("", appHandler.QuizHandler.CreateQuiz)
			quizzes.PUT("/:id", appHandler.QuizHandler.UpdateQuiz)
			quizzes.DELETE("/:id", appHandler.QuizHandler.DeleteQuiz)
			quizzes.POST("/submit", appHandler.QuizHandler.SubmitQuiz)
		}
	}
}
