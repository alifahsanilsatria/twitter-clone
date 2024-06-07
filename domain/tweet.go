package domain

import "context"

type UndoLikesRequestPayload struct {
	TweetId int32 `json:"tweet_id"`
}

type LikeTweetRequestPayload struct {
	TweetId int32 `json:"tweet_id"`
}

type UndoRetweetRequestPayload struct {
	TweetId int32 `json:"tweet_id"`
}

type RetweetRequestPayload struct {
	TweetId int32 `json:"tweet_id"`
}
type PublishTweetRequestPayload struct {
	ParentId int32  `json:"parent_id"`
	Content  string `json:"content"`
}

type DeleteTweetRequestPayload struct {
	TweetId int32 `json:"tweet_id"`
	UserId  int32 `json:"user_id"`
}

type PublishTweetUsecaseParam struct {
	Token    string
	ParentId int32
	Content  string
}

type PublishTweetUsecaseResult struct {
	TweetId int32
}

type DeleteTweetUsecaseParam struct {
	Token   string
	TweetId int32
}

type DeleteTweetUsecaseResult struct {
	TweetId int32
}

type RetweetUsecaseParam struct {
	Token   string
	TweetId int32
}

type RetweetUsecaseResult struct{}

type UndoRetweetUsecaseParam struct {
	Token   string
	TweetId int32
}

type UndoRetweetUsecaseResult struct{}

type LikeTweetUsecaseParam struct {
	Token   string
	TweetId int32
}

type LikeTweetUsecaseResult struct{}

type UndoLikesUsecaseParam struct {
	Token   string
	TweetId int32
}

type UndoLikesUsecaseResult struct{}

type SeeTweetDetailsUsecaseParam struct {
	Token   string
	TweetId int32
}

type SeeTweetDetailsUsecaseResult struct {
	TweetDetails []SeeTweetDetailsUsecaseResult_TweetDetails
}

type SeeTweetDetailsUsecaseResult_TweetDetails struct {
	TweetId      int32
	Username     string
	CompleteName string
	Content      string
	CountRetweet int32
	CountLikes   int32
	CountReplies int32

	IsParentTweet  bool
	IsCurrentTweet bool
	IsChildTweet   bool
}

type GetListOfUserTimelineTweetsParam struct {
	Token string
}

type GetListOfUserTimelineTweetsResult struct {
	Tweets []GetListOfUserTimelineTweetsResult_TweetDetails
}

type GetListOfUserTimelineTweetsResult_TweetDetails struct {
	TweetId      int32
	Username     string
	CompleteName string
	Content      string
	CountRetweet int32
	CountLikes   int32
	CountReplies int32
}

type GetListOfAUserTimelineTweetsParam struct {
	Token    string
	Username string
}

type GetListOfAUserTimelineTweetsResult struct {
	Tweets []GetListOfAUserTimelineTweetsResult_TweetDetails
}

type GetListOfAUserTimelineTweetsResult_TweetDetails struct {
	TweetId      int32
	Username     string
	CompleteName string
	Content      string
	CountRetweet int32
	CountLikes   int32
	CountReplies int32
}

type TweetUsecase interface {
	PublishTweet(ctx context.Context, param PublishTweetUsecaseParam) (PublishTweetUsecaseResult, error)
	DeleteTweet(ctx context.Context, param DeleteTweetUsecaseParam) (DeleteTweetUsecaseResult, error)
	Retweet(ctx context.Context, param RetweetUsecaseParam) (RetweetUsecaseResult, error)
	UndoRetweet(ctx context.Context, param UndoRetweetUsecaseParam) (UndoRetweetUsecaseResult, error)
	LikeTweet(ctx context.Context, param LikeTweetUsecaseParam) (LikeTweetUsecaseResult, error)
	UndoLikes(ctx context.Context, param UndoLikesUsecaseParam) (UndoLikesUsecaseResult, error)
	SeeTweetDetails(ctx context.Context, param SeeTweetDetailsUsecaseParam) (SeeTweetDetailsUsecaseResult, error)
	GetListOfUserTimelineTweets(ctx context.Context, param GetListOfUserTimelineTweetsParam) (GetListOfUserTimelineTweetsResult, error)
	GetListOfAUserTimelineTweets(ctx context.Context, param GetListOfAUserTimelineTweetsParam) (GetListOfAUserTimelineTweetsResult, error)
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

type UpsertRetweetParam struct {
	UserId  int32
	TweetId int32
}

type UpsertRetweetResult struct{}

type DeleteRetweetParam struct {
	UserId  int32
	TweetId int32
}

type DeleteRetweetResult struct{}

type UpsertLikesParam struct {
	UserId  int32
	TweetId int32
}

type UpsertLikesResult struct{}

type DeleteLikesParam struct {
	UserId  int32
	TweetId int32
}

type DeleteLikesResult struct{}

type GetChildrenDataByTweetIdParam struct {
	TweetId int32
}

type GetChildrenDataByTweetIdResult struct {
	ChildTweet []GetChildrenDataByTweetIdResult_ChildTweet
}

type GetChildrenDataByTweetIdResult_ChildTweet struct {
	TweetId      int32
	Username     string
	CompleteName string
	Content      string
	CountRetweet int32
	CountLikes   int32
	CountReplies int32
}

type GetParentsDataByTweetIdParam struct {
	TweetId int32
}

type GetParentsDataByTweetIdResult struct {
	Parent []GetParentsDataByTweetIdResult_Parent
}

type GetParentsDataByTweetIdResult_Parent struct {
	TweetId      int32
	Username     string
	CompleteName string
	Content      string
	CountRetweet int32
	CountLikes   int32
	CountReplies int32
}

type GetTweetByIdParam struct {
	TweetId int32
}

type GetTweetByIdResult struct {
	TweetId      int32
	Username     string
	CompleteName string
	Content      string
	CountRetweet int32
	CountLikes   int32
	CountReplies int32
}

type GetListOfUserFollowingTweetsParam struct {
	UserId int32
}

type GetListOfUserFollowingTweetsResult struct {
	Tweets []GetListOfUserFollowingTweetsResult_Tweet
}

type GetListOfUserFollowingTweetsResult_Tweet struct {
	TweetId      int32
	Username     string
	CompleteName string
	Content      string
	CountRetweet int32
	CountLikes   int32
	CountReplies int32
}

type GetListOfAUserTweetsParam struct {
	UserId int32
}

type GetListOfAUserTweetsResult struct {
	Tweets []GetListOfAUserTweetsResult_Tweet
}

type GetListOfAUserTweetsResult_Tweet struct {
	TweetId      int32
	Username     string
	CompleteName string
	Content      string
	CountRetweet int32
	CountLikes   int32
	CountReplies int32
}

type TweetRepository interface {
	CreateNewTweet(ctx context.Context, param CreateNewTweetParam) (CreateNewTweetResult, error)
	GetTweetByIdAndUserId(ctx context.Context, param GetTweetByIdAndUserIdParam) (GetTweetByIdAndUserIdResult, error)
	DeleteTweetById(ctx context.Context, param DeleteTweetByIdParam) (DeleteTweetByIdResult, error)
	UpsertRetweet(ctx context.Context, param UpsertRetweetParam) (UpsertRetweetResult, error)
	DeleteRetweet(ctx context.Context, param DeleteRetweetParam) (DeleteRetweetResult, error)
	UpsertLikes(ctx context.Context, param UpsertLikesParam) (UpsertLikesResult, error)
	DeleteLikes(ctx context.Context, param DeleteLikesParam) (DeleteLikesResult, error)
	GetChildrenDataByTweetId(ctx context.Context, param GetChildrenDataByTweetIdParam) (GetChildrenDataByTweetIdResult, error)
	GetParentsDataByTweetId(ctx context.Context, param GetParentsDataByTweetIdParam) (GetParentsDataByTweetIdResult, error)
	GetTweetById(ctx context.Context, param GetTweetByIdParam) (GetTweetByIdResult, error)
	GetListOfUserFollowingTweets(ctx context.Context, param GetListOfUserFollowingTweetsParam) (GetListOfUserFollowingTweetsResult, error)
	GetListOfAUserTweets(ctx context.Context, param GetListOfAUserTweetsParam) (GetListOfAUserTweetsResult, error)
}
