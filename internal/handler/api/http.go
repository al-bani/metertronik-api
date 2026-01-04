package api

import (
	"metertronik/internal/service/http"
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

	var data *service.DailyActivityResponse
	var err error

	if date == "" {
		data, err = h.apiService.DayNowActivity(c.Request.Context(), deviceID)
	} else {
		data, err = h.apiService.DailyActivity(c.Request.Context(), deviceID, date)
	}

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

	data, err := h.apiService.DailyRange(c.Request.Context(), id, startDate, endDate, lastDate, 10)

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

func (h *ApiHandler) GetMonthlyList(c *gin.Context) {
	id := c.Param("id")
	date := c.Query("date")

	data, err := h.apiService.MonthlyList(c.Request.Context(), id, date)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"id":      id,
		"data":    data,
	})
}

func (h *ApiHandler) UserPairing(c *gin.Context) {
	var reqUserPairing struct {
		DeviceId string `json:"device_id" binding:"required"`
	}
	c.ShouldBindJSON(&reqUserPairing)

	userId := c.GetInt64("user_id")

	pairing, err := h.apiService.UserPairing(c.Request.Context(), reqUserPairing.DeviceId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"message": "Failed to pair user",

		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"pairing_token": pairing,
	})
}

func (h *ApiHandler) DevicePairing(c *gin.Context) {
	var reqDevicePairing struct {
		DeviceId string `json:"device_id" binding:"required"`
		DeviceSecret string `json:"device_secret" binding:"required"`
		PairToken string `json:"pair_token" binding:"required"`
	}
	c.ShouldBindJSON(&reqDevicePairing)

	err := h.apiService.DevicePairing(c.Request.Context(), reqDevicePairing.DeviceId, reqDevicePairing.DeviceSecret, reqDevicePairing.PairToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"message": "Failed to pair device",
			"is_paired": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"is_paired": true,
	})
}

func (h *ApiHandler) PairingStatus(c *gin.Context) {
	deviceId := c.Param("id")

	pairing, err := h.apiService.PairingStatus(c.Request.Context(), deviceId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"is_paired": pairing,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_paired": pairing,
	})
}