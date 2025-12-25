package http

import (
	handler "metertronik/internal/handler/api"
	"metertronik/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, apiHandler *handler.ApiHandler, authHandler *handler.AuthHandler) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/verify", authHandler.Verify)
	}

	api := r.Group("/v1/api")
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
