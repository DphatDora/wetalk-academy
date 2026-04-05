package router

import (
	"strings"
	"wetalk-academy/config"
	repository "wetalk-academy/internal/infrastructure/db/repository"
	"wetalk-academy/internal/infrastructure/judge0"
	"wetalk-academy/internal/interface/handler"
	"wetalk-academy/internal/interface/middleware"
	"wetalk-academy/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppHandler struct {
	topicHandler   *handler.TopicHandler
	lessonHandler  *handler.LessonHandler
	contentHandler *handler.ContentHandler
	judge0Handler  *handler.Judge0Handler
	quizHandler    *handler.QuizHandler
}

func SetupRoutes(mongoDB *mongo.Database, conf *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware(conf.App.Whitelist))

	// repositories
	topicRepo := repository.NewTopicRepository(mongoDB)
	lessonRepo := repository.NewLessonRepository(mongoDB)
	contentRepo := repository.NewContentRepository(mongoDB)
	quizRepo := repository.NewQuizRepository(mongoDB)
	quizSubmissionRepo := repository.NewQuizSubmissionRepository(mongoDB)

	// clients
	judge0Client := judge0.NewClient(conf)

	// services
	topicService := service.NewTopicService(topicRepo, lessonRepo)
	lessonService := service.NewLessonService(lessonRepo, topicRepo, contentRepo)
	contentService := service.NewContentService(contentRepo, lessonRepo, topicRepo)
	judge0Service := service.NewJudge0Service(judge0Client)
	quizService := service.NewQuizService(quizRepo, quizSubmissionRepo, lessonRepo, topicRepo)

	// handlers
	appHandler := &AppHandler{
		topicHandler:   handler.NewTopicHandler(topicService),
		lessonHandler:  handler.NewLessonHandler(lessonService),
		contentHandler: handler.NewContentHandler(contentService),
		judge0Handler:  handler.NewJudge0Handler(judge0Service),
		quizHandler:    handler.NewQuizHandler(quizService),
	}

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

func setupPublicRoutes(rg *gin.RouterGroup, conf *config.Config, appHandler *AppHandler) {
	rg.Use(middleware.OptionalAuthMiddleware(conf))

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
		lessons.GET("/:slug/content", appHandler.contentHandler.GetContent)
		lessons.GET("/:slug/quiz", appHandler.quizHandler.GetQuizzesByLessonSlug)
	}

	quizzes := rg.Group("/quizzes")
	{
		quizzes.GET("/:id", appHandler.quizHandler.GetQuizByID)
	}

	judge0 := rg.Group("/judge0")
	{
		judge0.POST("/submit", appHandler.judge0Handler.SubmitCode)
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

			lessons.POST("/:slug/content", appHandler.contentHandler.CreateContent)
			lessons.PUT("/:slug/content", appHandler.contentHandler.UpdateContent)
			lessons.DELETE("/:slug/content", appHandler.contentHandler.DeleteContent)
		}

		quizzes := protected.Group("/quizzes")
		{
			quizzes.POST("", appHandler.quizHandler.CreateQuiz)
			quizzes.PUT("/:id", appHandler.quizHandler.UpdateQuiz)
			quizzes.DELETE("/:id", appHandler.quizHandler.DeleteQuiz)
			quizzes.POST("/submit", appHandler.quizHandler.SubmitQuiz)
		}
	}
}
