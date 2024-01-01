package domain

import (
	"context"
)

type SignUpParam struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	CompleteName string `json:"complete_name"`
}

type SignUpResult struct {
	Id int64 `json:"id"`
}

type SignInParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResult struct {
	Token string `json:"token"`
}

type SignOutParam struct {
	Token string
}

type SignOutResult struct {
}

type UserUsecase interface {
	SignUp(ctx context.Context, param SignUpParam) (SignUpResult, error)
	SignIn(ctx context.Context, param SignInParam) (SignInResult, error)
	SignOut(ctx context.Context, param SignOutParam) (SignOutResult, error)
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

type UserRepository interface {
	GetUserByUsernameOrEmail(ctx context.Context, param GetUserByUsernameOrEmailParam) (GetUserByUsernameOrEmailResult, error)
	CreateNewUserAccount(ctx context.Context, param CreateNewUserAccountParam) (CreateNewUserAccountResult, error)
	GetUserByUsername(ctx context.Context, param GetUserByUsernameParam) (GetUserByUsernameResult, error)
}
