package controllers

import (
	"github.com/gin-gonic/gin"
	"loan/usecase"
	"net/http"
	"strconv"

)

type LogController struct {
	logUsecase *usecase.LogUsecase
}

func NewLogController(lu *usecase.LogUsecase) *LogController {
	return &LogController{
		logUsecase: lu,
	}
}

func (lc *LogController) ViewLogs(c *gin.Context) {
	// Role management and admin ID retrieval
	Role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: role not found"})
		return
	}

	if Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: admin ID not found"})
		return
	}

	// Query parameters
	event := c.Query("event")
	order := c.Query("order")

	const defaultLimit, defaultOffset = 10, 0

	limit, err := strconv.ParseInt(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)), 10, 64)
	if err != nil {
		limit = defaultLimit
	}

	offset, err := strconv.ParseInt(c.DefaultQuery("offset", strconv.Itoa(defaultOffset)), 10, 64)
	if err != nil {
		offset = defaultOffset
	}

	// Call the usecase to retrieve logs
	logs, total, uerr := lc.logUsecase.ViewLogs(event, order, int(limit), int(offset), adminID.(string))
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{"error": uerr.Message})
		return
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	// Return the logs along with pagination information
	c.JSON(http.StatusOK, gin.H{
		"logs":        logs,
		"current_page": offset,
		"per_page":    limit,
		"total":       total,
		"total_pages": totalPages,
	})
}
