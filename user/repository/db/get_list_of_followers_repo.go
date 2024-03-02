package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (repo *userRepository) GetListOfFollowersRepo(ctx context.Context, param domain.GetListOfFollowersRepoParam) (domain.GetListOfFollowersRepoResult, error) {
	ctx, span := repo.tracer.Start(ctx, "repository.GetListOfFollowersRepo", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "userRepository.GetListOfFollowersRepo",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	query := `
		select u.id, u.username, u.complete_name
		from user_following uf join user u
		on uf.user_id = u.id
		where uf.following_user_id = $1
		and u.is_deleted = false
		and uf.is_deleted = false
	`

	args := []interface{}{
		param.UserId,
	}

	result := domain.GetListOfFollowersRepoResult{
		Users: []domain.GetListOfFollowersRepoResult_User{},
	}

	logData["query"] = query
	logData["args"] = fmt.Sprintf("%+v", args)

	queryContextResp, errQueryContext := repo.db.QueryContext(ctx, query, args...)
	if errQueryContext != nil {
		logData["error_query_context"] = errQueryContext.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQueryContext).
			Errorln("error on querycontext")
		span.End()
		return result, errQueryContext
	}

	for queryContextResp.Next() {
		user := domain.GetListOfFollowersRepoResult_User{}
		errScan := queryContextResp.Scan(&user.UserId, &user.Username, &user.CompleteName)
		if errScan != nil && errScan != sql.ErrNoRows {
			logData["error_scan"] = errScan.Error()
			repo.logger.
				WithFields(logData).
				WithError(errScan).
				Errorln("error on scan")
			span.End()
			return result, errScan
		}
		result.Users = append(result.Users, user)
	}

	logData["response"] = fmt.Sprintf("%+v", result)
	repo.logger.
		WithFields(logData).
		Infoln("success GetListOfFollowersRepo")
	span.End()

	return result, nil
}
