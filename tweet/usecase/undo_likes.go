package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (uc *tweetUsecase) UndoLikes(ctx context.Context, param domain.UndoLikesUsecaseParam) (domain.UndoLikesUsecaseResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.UndoLikes", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetUsecase.UndoLikes",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	getUserSessionByTokenParam := domain.GetUserSessionByTokenParam{
		Token: param.Token,
	}

	logData["get_user_session_by_token_param"] = fmt.Sprintf("%+v", getUserSessionByTokenParam)

	userSession, errGetUserSession := uc.userSessionRepository.GetUserSessionByToken(ctx, getUserSessionByTokenParam)
	if errGetUserSession != nil {
		logData["error_get_user_session"] = errGetUserSession.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetUserSession).
			Errorln("error on GetUserSessionByToken")
		span.End()
		return domain.UndoLikesUsecaseResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.UndoLikesUsecaseResult{}, errors.New("invalid or expired token")
	}

	deleteLikesParam := domain.DeleteLikesParam{
		UserId:  userSession.UserId,
		TweetId: param.TweetId,
	}

	deleteLikesResult, errDeleteLikes := uc.tweetRepository.DeleteLikes(ctx, deleteLikesParam)
	if errDeleteLikes != nil {
		logData["error_delete_likes"] = errDeleteLikes.Error()
		uc.logger.
			WithFields(logData).
			WithError(errDeleteLikes).
			Errorln("error on DeleteLikes")
		span.End()
		return domain.UndoLikesUsecaseResult{}, errDeleteLikes
	}

	logData["delete_likes_result"] = fmt.Sprintf("%+v", deleteLikesResult)
	uc.logger.
		WithFields(logData).
		Infoln("success undo likes")

	span.End()

	result := domain.UndoLikesUsecaseResult{}
	return result, nil
}
