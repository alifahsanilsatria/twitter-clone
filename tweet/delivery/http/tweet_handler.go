package http

import (
	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	trace "go.opentelemetry.io/otel/trace"
)

type tweetHandler struct {
	tweetUsecase domain.TweetUsecase
	logger       *logrus.Logger
	tracer       trace.Tracer
}

func NewTweetHandler(
	e *echo.Echo,
	us domain.TweetUsecase,
	logger *logrus.Logger,
	tracer trace.Tracer,
) {
	handler := &tweetHandler{
		tweetUsecase: us,
		logger:       logger,
		tracer:       tracer,
	}
	e.POST("/tweet", handler.PublishTweet)
	e.DELETE("/tweet", handler.DeleteTweet)
	e.POST("/retweet", handler.Retweet)
	e.DELETE("/retweet", handler.UndoRetweet)
	e.POST("/likes", handler.LikeTweet)
	e.DELETE("/likes", handler.UndoLikes)
}
