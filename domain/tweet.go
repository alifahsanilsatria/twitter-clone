package domain

import "context"

type PublishTweetRequestPayload struct {
	ParentId int32  `json:"parent_id"`
	Content  string `json:"content"`
}

type DeleteTweetRequestPayload struct {
	TweetId int32 `json:"tweet_id"`
	UserId  int32 `json:"user_id"`
}

type PublishTweetParam struct {
	Token    string
	ParentId int32
	Content  string
}

type PublishTweetResult struct {
	TweetId int32
}

type DeleteTweetParam struct {
	Token   string
	TweetId int32
	UserId  int32
}

type DeleteTweetResult struct {
	TweetId int32
}

type TweetUsecase interface {
	PublishTweet(ctx context.Context, param PublishTweetParam) (PublishTweetResult, error)
	DeleteTweet(ctx context.Context, param DeleteTweetParam) (DeleteTweetResult, error)
}

type CreateNewTweetParam struct {
	UserId   int32
	ParentId int32
	Content  string
}

type CreateNewTweetResult struct {
	TweetId int32
}

type GetTweetByIdAndUserIdParam struct {
	TweetId int32
	UserId  int32
}

type GetTweetByIdAndUserIdResult struct {
	TweetId int32
}

type DeleteTweetByIdParam struct {
	TweetId int32
}

type DeleteTweetByIdResult struct {
	TweetId int32
}

type TweetRepository interface {
	CreateNewTweet(ctx context.Context, param CreateNewTweetParam) (CreateNewTweetResult, error)
	GetTweetByIdAndUserId(ctx context.Context, param GetTweetByIdAndUserIdParam) (GetTweetByIdAndUserIdResult, error)
	DeleteTweetById(ctx context.Context, param DeleteTweetByIdParam) (DeleteTweetByIdResult, error)
}
