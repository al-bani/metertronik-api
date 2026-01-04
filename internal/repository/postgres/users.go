package postgres

import (
	"context"
	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"
	"errors"

	"fmt"
	"gorm.io/gorm"
)

type UsersRepoPostgres struct {
	db *gorm.DB
}

func NewUsersRepoPostgres(db *gorm.DB) repository.UsersRepoPostgres {
	return &UsersRepoPostgres{
		db: db,
	}
}

func (r *UsersRepoPostgres) CreateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Table("users").Create(user).Error
}

func (r *UsersRepoPostgres) GetUser(ctx context.Context, email string, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Table("users").Where("email = ? OR username = ?", email, username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepoPostgres) UpdateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Table("users").Where("email = ?", user.Email).Updates(user).Error
}

func (r *UsersRepoPostgres) CheckDeviceUser(ctx context.Context, userId int64) (bool, error) {
	var deviceUser entity.DeviceUser

	result := r.db.WithContext(ctx).
		Table("device_users").
		Where("user_id = ?", userId).
		First(&deviceUser)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if result.Error != nil {
		return false, fmt.Errorf("failed to check device user: %w", result.Error)
	}

	return true, nil
}