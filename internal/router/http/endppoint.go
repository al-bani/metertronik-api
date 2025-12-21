package http

import (
	handler "metertronik/internal/handler/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, handler *handler.ApiHandler) {
	api := r.Group("/v1/api")

	{
		api.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "OK",
				"version": "1.0.0",
			})
		})
		api.GET("/daily/:id", handler.GetDailyList)
		api.GET("/daily/:id/detail", handler.GetSpecificDailyActivity)
		api.GET("/daily/:id/range", handler.GetDailyRange)
		api.GET("/monthly/:id", handler.GetMonthlyList)

		// api.GET("/daily/summary", func(ctx *gin.Context) {

		// })
	}

}
