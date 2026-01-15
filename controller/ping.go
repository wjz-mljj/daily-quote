package controller

import (
	"a-sentence/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	msg := service.PingService()
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": msg,
	})
}
