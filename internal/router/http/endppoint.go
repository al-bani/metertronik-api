package http

import (
	handler "metertronik/internal/handler/api"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, handler *handler.ApiHandler) {
	api := r.Group("/v1/api")

	{
		api.GET("/daily/:id", handler.GetDailyList)
		api.GET("/daily/:id/detail", handler.GetSpecificDailyActivity)
		api.GET("/daily/:id/range", handler.GetDailyRange)

		// api.GET("/daily/summary", func(ctx *gin.Context) {

		// })
	}

}
