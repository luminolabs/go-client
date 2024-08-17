package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartAPIServer(port string) error {
	router := gin.Default()

	// Define your API endpoints here
	router.GET("/status", getStatus)

	return router.Run(":" + port)
}

func getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "running",
	})
}
