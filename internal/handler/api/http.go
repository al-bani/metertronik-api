package api

import (
	"metertronik/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	apiService *service.ApiService
}

func NewApiHandler(apiService *service.ApiService) *ApiHandler {
	return &ApiHandler{
		apiService: apiService,
	}
}

func (h *ApiHandler) GetSpecificDailyActivity(c *gin.Context) {
	deviceID := c.Param("id")
	date := c.Query("date")

	data, err := h.apiService.DailyActivity(c.Request.Context(), deviceID, date)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "OK",
		"device_id": deviceID,
		"date":      date,
		"data":      data,
	})
}

func (h *ApiHandler) GetDailyList(c *gin.Context) {
	id := c.Param("id")
	time := c.Query("time")
	tariff := c.Query("tariff")

	data, err := h.apiService.DailyList(c.Request.Context(), id, time, tariff)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"id":      id,
		"data":    data,
	})
}

func (h *ApiHandler) GetDailyRange(c *gin.Context) {
	id := c.Param("id")
	startDate := c.Query("start")
	endDate := c.Query("end")

	data, err := h.apiService.DailyRange(c.Request.Context(), id, startDate, endDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"id":      id,
		"start":   startDate,
		"end":     endDate,
		"data":    data,
	})
}
