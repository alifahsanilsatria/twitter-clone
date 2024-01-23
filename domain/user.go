package domain

import (
	"context"
)

type UnfollowRequestPayload struct {
	FollowingUserId int32 `json:"following_user_id"`
}
type FollowUserRequestPayload struct {
	FollowingUserId int32 `json:"following_user_id"`
}

type SignUpRequestPayload struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	CompleteName string `json:"complete_name"`
}

type SignUpUsecaseParam struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	CompleteName string `json:"complete_name"`
}

type SignUpResult struct {
	Id int64 `json:"id"`
}

type SignInRequestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInUsecaseParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResult struct {
	Token string `json:"token"`
}

type SignOutUsecaseParam struct {
	Token string
}

type SignOutResult struct {
}

type SeeProfileDetailsParam struct {
	Token string
}

type SeeProfileDetailsResult struct {
	Username     string
	Email        string
	CompleteName string
	CreatedAt    string
}

type FollowUserParam struct {
	Token           string
	FollowingUserId int32
}

type FollowUserResult struct {
	Id int32
}

type UnfollowUserParam struct {
	Token           string
	FollowingUserId int32
}

type UnfollowUserResult struct {
	Id int32
}

type UserUsecase interface {
	SignUp(ctx context.Context, param SignUpUsecaseParam) (SignUpResult, error)
	SignIn(ctx context.Context, param SignInUsecaseParam) (SignInResult, error)
	SignOut(ctx context.Context, param SignOutUsecaseParam) (SignOutResult, error)
	SeeProfileDetails(ctx context.Context, param SeeProfileDetailsParam) (SeeProfileDetailsResult, error)
	FollowUser(ctx context.Context, param FollowUserParam) (FollowUserResult, error)
	UnfollowUser(ctx context.Context, param UnfollowUserParam) (UnfollowUserResult, error)
}

type GetUserByUsernameOrEmailParam struct {
	Username string
	Email    string
}

type GetUserByUsernameOrEmailResult struct {
	Id int64
}

type CreateNewUserAccountParam struct {
	Username       string
	HashedPassword string
	Email          string
	CompleteName   string
}

type CreateNewUserAccountResult struct {
	Id int64
}

type GetUserByUsernameParam struct {
	Username string
}

type GetUserByUsernameResult struct {
	Id             int32
	HashedPassword string
}

type GetUserByUserIdParam struct {
	UserId int32
}

type GetUserByUserIdResult struct {
	UserId       int32
	Username     string
	Email        string
	CompleteName string
	CreatedAt    string
}

type UpsertUserFollowingParam struct {
	UserId          int32
	FollowingUserId int32
}

type UpsertUserFollowingResult struct {
	Id int32
}

type DeleteUserFollowingParam struct {
	UserId          int32
	FollowingUserId int32
}

type DeleteUserFollowingResult struct {
	Id int32
}

type UserRepository interface {
	GetUserByUsernameOrEmail(ctx context.Context, param GetUserByUsernameOrEmailParam) (GetUserByUsernameOrEmailResult, error)
	CreateNewUserAccount(ctx context.Context, param CreateNewUserAccountParam) (CreateNewUserAccountResult, error)
	GetUserByUsername(ctx context.Context, param GetUserByUsernameParam) (GetUserByUsernameResult, error)
	GetUserByUserId(ctx context.Context, param GetUserByUserIdParam) (GetUserByUserIdResult, error)
	UpsertUserFollowing(ctx context.Context, param UpsertUserFollowingParam) (UpsertUserFollowingResult, error)
	DeleteUserFollowing(ctx context.Context, param DeleteUserFollowingParam) (DeleteUserFollowingResult, error)
}
