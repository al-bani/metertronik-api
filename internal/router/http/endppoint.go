package http

import (
	handler "metertronik/internal/handler/api"
	"metertronik/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, apiHandler *handler.ApiHandler, authHandler *handler.AuthHandler) {
	api := r.Group("/v1/api")

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.Refresh)
		auth.GET("/logout/:id", authHandler.Logout)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/verify-otp", authHandler.VerifyOtp)
		auth.POST("/resend-otp", authHandler.ResendOtp)
		auth.POST("/request-reset-password", authHandler.RequestResetPassword)
	}


	api.Use(middleware.JWTMiddleware())
	{
		api.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "OK",
				"version": "1.0.0",
			})
		})
		api.GET("/daily/:id", apiHandler.GetDailyList)
		api.GET("/daily/:id/detail", apiHandler.GetSpecificDailyActivity)
		api.GET("/daily/:id/range", apiHandler.GetDailyRange)
		api.GET("/monthly/:id", apiHandler.GetMonthlyList)

		// api.GET("/daily/summary", func(ctx *gin.Context) {

		// })
	}

}
