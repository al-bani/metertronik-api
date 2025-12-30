package repository

import (
	"context"
	"metertronik/internal/domain/entity"
)

type UsersRepoPostgres interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUser(ctx context.Context, email string, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
}

type UsersRepoRedis interface {
	ResetExpiredToken(ctx context.Context, identifier int64, token string) error
	SetVerificationCodeOtp(ctx context.Context, identifier string, hashCode string) error
	GetVerificationCodeOtp(ctx context.Context, identifier string) (string, error)
	SetToken(ctx context.Context, identifier int64, token string) error
	DeleteToken(ctx context.Context, token string) error
}