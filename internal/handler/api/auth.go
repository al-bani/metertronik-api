package api

import (
	"log"
	"metertronik/internal/domain/entity"
	service "metertronik/internal/service/http"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RequestUser struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RequestUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	if req.Email == "" || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Email, username, dan password tidak boleh kosong",
		})
		return
	}

	user := &entity.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Role:     "user",
		Status:   "active",
		Verified: false,
	}

	err := h.authService.RegisterController(c.Request.Context(), user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Registration failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"data": gin.H{
			"email":    user.Email,
			"username": user.Username,
			"role":     user.Role,
			"status":   user.Status,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" && req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Email atau username harus diisi (minimal salah satu)",
		})
		return
	}

	user := &entity.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	tokenResponse, err := h.authService.LoginController(c.Request.Context(), user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Login failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    tokenResponse,
	})

}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
		UserID       int64  `json:"user_id" binding:"required"`
	}
	c.ShouldBindJSON(&req)

	tokenResponse, err := h.authService.RefreshController(c.Request.Context(), req.UserID, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Refresh failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Refresh successful",
		"data":    tokenResponse,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	log.Println("Logout Handler")
	var reqLogout struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqLogout); err != nil {
		log.Println("Logout Handler 1")
		c.JSON(http.StatusBadRequest, gin.H{
			"status_logout": false,
			"error":         "Invalid request",
			"message":       err.Error(),
		})
		return
	}
	log.Println("Logout 1")
	err := h.authService.LogoutController(c.Request.Context(), reqLogout.RefreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_logout": false,
			"error":         "Logout failed",
			"message":       err.Error(),
		})
		return
	}
	log.Println("Logout 2")
	c.JSON(http.StatusOK, gin.H{
		"status_logout": true,
		"message":       "Logout successful",
	})
}

func (h *AuthHandler) RequestResetPassword(c *gin.Context) {
	var reqRequestResetPassword struct {
		Email string `json:"email" binding:"required,email"`
	}
	c.ShouldBindJSON(&reqRequestResetPassword)

	err := h.authService.RequestResetPasswordController(c.Request.Context(), reqRequestResetPassword.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Request reset password failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Request reset password sent to email",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var reqResetPassword struct {
		Email    string `json:"email" binding:"required,email"`
		Otp      string `json:"otp" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	c.ShouldBindJSON(&reqResetPassword)

	err := h.authService.ResetPasswordController(c.Request.Context(), reqResetPassword.Email, reqResetPassword.Otp, reqResetPassword.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Reset password failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reset password successful",
	})
}

func (h *AuthHandler) VerifyOtp(c *gin.Context) {
	var reqVerifyOtp struct {
		Email string `json:"email" binding:"required,email"`
		Otp   string `json:"otp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqVerifyOtp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_verify": false,
			"error":         "Invalid request",
			"message":       err.Error(),
		})
		return
	}

	err := h.authService.VerifyOtpController(c.Request.Context(), reqVerifyOtp.Email, reqVerifyOtp.Otp)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_verify": false,
			"error":         "Verify OTP failed",
			"message":       err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_verify": true,
		"message":       "Verify OTP successful",
	})
}

func (h *AuthHandler) ResendOtp(c *gin.Context) {
	var reqResendOtp struct {
		Email string `json:"email" binding:"required,email"`
	}
	c.ShouldBindJSON(&reqResendOtp)

	err := h.authService.ResendOtpController(c.Request.Context(), reqResendOtp.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_resend": false,
			"error":         "Resend OTP failed",
			"message":       err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_resend": true,
		"message":       "Resend OTP successful",
	})
}

func (h *AuthHandler) CheckId(c *gin.Context) {
	log.Println("CheckId Handler")
	var reqCheckId struct {
		Id string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqCheckId); err != nil {
		log.Println("CheckId Handler 1")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}
	log.Println("CheckId Handler 2")
	avail, err := h.authService.CheckIdController(c.Request.Context(), reqCheckId.Id)
	if err != nil {
		// Validator di service mengembalikan error dengan prefix "invalid id:"
		if strings.HasPrefix(err.Error(), "invalid id:") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Check ID failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available": avail,
		"message":   "Check ID successful",
	})
}
