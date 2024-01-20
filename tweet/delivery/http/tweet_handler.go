package http

import (
	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

type tweetHandler struct {
	tweetUsecase domain.TweetUsecase
	logger       *logrus.Logger
}

func NewTweetHandler(
	e *echo.Echo,
	us domain.TweetUsecase,
	logger *logrus.Logger,
) {
	handler := &tweetHandler{
		tweetUsecase: us,
		logger:       logger,
	}
	e.POST("/tweet", handler.PublishTweet)
	e.DELETE("/tweet", handler.DeleteTweet)
	e.POST("/retweet", handler.Retweet)
	e.DELETE("/retweet", handler.UndoRetweet)
	e.POST("/likes", handler.LikeTweet)
	e.DELETE("/likes", handler.UndoLikes)
}
