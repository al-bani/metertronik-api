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
	lastDate := c.Query("last")

	data, err := h.apiService.DailyList(c.Request.Context(), id, time, tariff, lastDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var lastDateData string
	if data != nil && len(*data) > 0 {
		lastDateData = (*data)[len(*data)-1].Day.Format()
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "OK",
		"id":        id,
		"data":      data,
		"last_date": lastDateData,
	})
}

func (h *ApiHandler) GetDailyRange(c *gin.Context) {
	id := c.Param("id")
	startDate := c.Query("start")
	endDate := c.Query("end")
	lastDate := c.Query("last")

	data, err := h.apiService.DailyRange(c.Request.Context(), id, startDate, endDate, lastDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var lastDateData string
	if data != nil && len(*data) > 0 {
		lastDateData = (*data)[len(*data)-1].Day.Format()
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "OK",
		"id":        id,
		"start":     startDate,
		"end":       endDate,
		"data":      data,
		"last_date": lastDateData,
	})
}
