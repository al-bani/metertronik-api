package service

import (
	"context"
	"errors"
	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"metertronik/pkg/utils/token"
)

type AuthService struct {
	postgresRepo  repository.UsersRepoPostgres
	redisAuthRepo repository.UsersRepoRedis
}

func NewAuthService(postgresRepo repository.UsersRepoPostgres, redisAuthRepo repository.UsersRepoRedis) *AuthService {
	return &AuthService{
		postgresRepo:  postgresRepo,
		redisAuthRepo: redisAuthRepo,
	}
}

type TokenResponse struct {
	User         *entity.User `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}


func (s *AuthService) RegisterController(ctx context.Context, user *entity.User) error {

	existingUser, err := s.postgresRepo.GetUser(ctx, user.Email, user.Username)
	if err == nil && existingUser != nil {
		return errors.New("user already exists, Email or Username already registered")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("failed to check existing user, " + err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}

	user.Password = string(hashedPassword)

	if user.Role == "" {
		user.Role = "user"
	}

	if user.Status == "" {
		user.Status = "active"
	}

	if err := s.postgresRepo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("user already exists, Email or Username already registered")
		}
		return errors.New("failed to create user, " + err.Error())
	}

	return nil
}

func (s *AuthService) LoginController(ctx context.Context, user *entity.User) (*TokenResponse, error) {
	existingUser, err := s.postgresRepo.GetUser(ctx, user.Email, user.Username)

	if err != nil {
		return nil, errors.New("failed to get user, " + err.Error())
	}

	if existingUser == nil {
		return nil, errors.New("user not found, Check your email, username or password")
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.New("user not found, Check your email, username or password")
	}

	refreshToken := token.GenerateRefreshToken()
	accessToken := token.GenerateAccessToken(existingUser.ID)

	s.redisAuthRepo.SetToken(ctx, existingUser.ID, refreshToken)

	return &TokenResponse{
		User:         existingUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *AuthService) RefreshController(ctx context.Context, userId int64, refreshToken string) (*TokenResponse, error) {
	newAccessToken := token.GenerateAccessToken(userId)

	err := s.redisAuthRepo.ResetExpired(ctx, userId, refreshToken)
	if err != nil {
		return nil, errors.New("failed to reset expired token, " + err.Error())
	}
	
	return &TokenResponse{
		User:         &entity.User{ID: userId},
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
	}, nil
}

// func (s *AuthService) LogoutController(ctx context.Context, user *entity.User) error {
// 	return s.postgresRepo.GetUser(ctx, user.Email)
// }

// func (s *AuthService) ResetPasswordController(ctx context.Context, user *entity.User) error {
// 	return s.postgresRepo.GetUser(ctx, user.Email)
// }

// func (s *AuthService) VerifyController(ctx context.Context, user *entity.User) error {
// 	return s.postgresRepo.GetUser(ctx, user.Email)
// }
