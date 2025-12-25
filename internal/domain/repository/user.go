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
	ResetExpired(ctx context.Context, identifier int64, token string) error
	TokenValidation(ctx context.Context, identifier string, token string) (bool, error)
	SetToken(ctx context.Context, identifier int64, token string) error
}