package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			time.Sleep(time.Second*5)
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}
	router.Run(":6000")
}
