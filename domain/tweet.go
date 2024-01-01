package domain

import "context"

type PublishTweetRequestPayload struct {
	ParentId int32  `json:"parent_id"`
	Content  string `json:"content"`
}

type PublishTweetParam struct {
	Token    string
	ParentId int32
	Content  string
}

type PublishTweetResult struct {
	TweetId int32
}

type TweetUsecase interface {
	PublishTweet(ctx context.Context, param PublishTweetParam) (PublishTweetResult, error)
}

type CreateNewTweetParam struct {
	UserId   int32
	ParentId int32
	Content  string
}

type CreateNewTweetResult struct {
	TweetId int32
}

type TweetRepository interface {
	CreateNewTweet(ctx context.Context, param CreateNewTweetParam) (CreateNewTweetResult, error)
}
