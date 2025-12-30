package service

import (
	"context"
	"errors"
	"metertronik/internal/domain/entity"
	"metertronik/internal/domain/repository"
	"metertronik/pkg/validator"

	"log"
	"metertronik/internal/handler/verification"
	"metertronik/pkg/utils"
	"metertronik/pkg/utils/token"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
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

	codeOtp, err := utils.GenerateOTP()

	hashCodeOtp, err := bcrypt.GenerateFromPassword([]byte(codeOtp), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash OTP code: " + err.Error())
	}

	s.redisAuthRepo.SetVerificationCodeOtp(ctx, user.Email, string(hashCodeOtp))

	err = verification.SendVerificationEmail(user.Email, codeOtp)
	if err != nil {
		return errors.New("failed to send verification email, " + err.Error())
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
	hashedRefreshToken := utils.Hashing(refreshToken)
	accessToken := token.GenerateAccessToken(existingUser.ID)
	log.Println("refresh token, LoginController : ", refreshToken)

	log.Println("hashedRefreshToken, LoginController : ", hashedRefreshToken)

	s.redisAuthRepo.SetToken(ctx, existingUser.ID, hashedRefreshToken)

	return &TokenResponse{
		User:         existingUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *AuthService) RefreshController(ctx context.Context, userId int64, refreshToken string) (*TokenResponse, error) {
	newAccessToken := token.GenerateAccessToken(userId)
	log.Println("refresh token, RefreshController : ", refreshToken)

	hashedRefreshToken := utils.Hashing(refreshToken)

	log.Println("hashedRefreshToken, RefreshController : ", hashedRefreshToken)

	err := s.redisAuthRepo.ResetExpiredToken(ctx, userId, hashedRefreshToken)

	if err != nil {
		return nil, errors.New("failed to reset expired token, " + err.Error())
	}

	return &TokenResponse{
		User:         &entity.User{ID: userId},
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) LogoutController(ctx context.Context, refreshToken string) error {
	log.Println("refresh token, LogoutController : ", refreshToken)
	hashedRefreshToken := utils.Hashing(refreshToken)

	return s.redisAuthRepo.DeleteToken(ctx, hashedRefreshToken)
}

func (s *AuthService) RequestResetPasswordController(ctx context.Context, email string) error {
	existingUser, err := s.postgresRepo.GetUser(ctx, email, "")
	if err != nil {
		return errors.New("failed to get user, " + err.Error())
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	codeOtp, err := utils.GenerateOTP()
	if err != nil {
		return errors.New("failed to generate OTP code, " + err.Error())
	}

	hashCodeOtp, err := bcrypt.GenerateFromPassword([]byte(codeOtp), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash OTP code, " + err.Error())
	}

	err = verification.SendVerificationEmail(email, codeOtp)
	if err != nil {
		return errors.New("failed to send verification email, " + err.Error())
	}

	s.redisAuthRepo.SetVerificationCodeOtp(ctx, email, string(hashCodeOtp))

	return nil
}

func (s *AuthService) ResetPasswordController(ctx context.Context, email string, otp string, password string) error {
	existingUser, err := s.postgresRepo.GetUser(ctx, email, "")
	if err != nil {
		return errors.New("failed to get user, " + err.Error())
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	storedHashOtp, err := s.redisAuthRepo.GetVerificationCodeOtp(ctx, email)
	if err != nil {
		return errors.New("failed to get verification code OTP, " + err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashOtp), []byte(otp))
	if err != nil {
		return errors.New("invalid verification code OTP")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password, " + err.Error())
	}

	existingUser.Password = string(hashedPassword)

	err = s.postgresRepo.UpdateUser(ctx, existingUser)
	if err != nil {
		return errors.New("failed to update user, " + err.Error())
	}

	return nil
}

func (s *AuthService) VerifyOtpController(ctx context.Context, email string, otp string) error {
	storedHashOtp, err := s.redisAuthRepo.GetVerificationCodeOtp(ctx, email)
	if err != nil {
		return errors.New("failed to get verification code OTP, " + err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashOtp), []byte(otp))

	if err != nil {
		return errors.New("invalid verification code OTP")
	}

	user, err := s.postgresRepo.GetUser(ctx, email, "")

	user.Verified = true

	err = s.postgresRepo.UpdateUser(ctx, user)
	if err != nil {
		return errors.New("failed to update user, " + err.Error())
	}

	return nil
}

func (s *AuthService) ResendOtpController(ctx context.Context, email string) error {
	existingUser, err := s.postgresRepo.GetUser(ctx, email, "")
	if err != nil {
		return errors.New("failed to get user, " + err.Error())
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	codeOtp, err := utils.GenerateOTP()
	if err != nil {
		return errors.New("failed to generate OTP code, " + err.Error())
	}

	hashCodeOtp, err := bcrypt.GenerateFromPassword([]byte(codeOtp), bcrypt.DefaultCost)

	if err != nil {
		return errors.New("failed to hash OTP code, " + err.Error())
	}

	s.redisAuthRepo.SetVerificationCodeOtp(ctx, email, string(hashCodeOtp))

	err = verification.SendVerificationEmail(email, codeOtp)
	if err != nil {
		return errors.New("failed to send verification email, " + err.Error())
	}

	return nil
}

func (s *AuthService) CheckIdController(ctx context.Context, id string) (bool, error) {
	if err := validator.ValidateControllerID(id); err != nil {
		return false, err
	}

	existingUser, err := s.postgresRepo.GetUser(ctx, "", id)
	if err != nil {
		// ID belum dipakai (record tidak ada) => available
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, errors.New("failed to get user, " + err.Error())
	}

	// Berjaga-jaga: kalau repo suatu saat mengembalikan (nil, nil)
	if existingUser == nil {
		return true, nil
	}

	// ID sudah dipakai
	return false, nil
}
