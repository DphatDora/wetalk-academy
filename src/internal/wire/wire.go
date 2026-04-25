//go:build wireinject
// +build wireinject

package wire

import (
	"wetalk-academy/config"
	domainrepo "wetalk-academy/internal/domain/repository"
	dbrepository "wetalk-academy/internal/infrastructure/db/repository"
	"wetalk-academy/internal/infrastructure/judge0"
	"wetalk-academy/internal/interface/handler"
	"wetalk-academy/internal/service"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppHandler struct {
	TopicHandler   *handler.TopicHandler
	LessonHandler  *handler.LessonHandler
	ContentHandler *handler.ContentHandler
	Judge0Handler  *handler.Judge0Handler
	QuizHandler    *handler.QuizHandler
}

var RepositorySet = wire.NewSet(
	dbrepository.NewTopicRepository,
	dbrepository.NewLessonRepository,
	dbrepository.NewContentRepository,
	dbrepository.NewQuizRepository,
	dbrepository.NewQuizSubmissionRepository,
	judge0.NewClient,
	wire.Bind(new(domainrepo.Judge0Repository), new(*judge0.Client)),
)

var ServiceSet = wire.NewSet(
	service.NewTopicService,
	service.NewLessonService,
	service.NewContentService,
	service.NewJudge0Service,
	service.NewQuizService,
)

var HandlerSet = wire.NewSet(
	handler.NewTopicHandler,
	handler.NewLessonHandler,
	handler.NewContentHandler,
	handler.NewJudge0Handler,
	handler.NewQuizHandler,
)

var ProviderSet = wire.NewSet(
	RepositorySet,
	ServiceSet,
	HandlerSet,
	wire.Struct(new(AppHandler), "*"),
)

func InitAppContainer(db *mongo.Database, conf *config.Config) *AppHandler {
	wire.Build(ProviderSet)
	return nil
}
