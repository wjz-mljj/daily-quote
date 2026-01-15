package middleware

import (
	"a-sentence/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic:", err)
				c.JSON(http.StatusOK, model.Response{
					Code:    5000,
					Message: "服务器内部错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
