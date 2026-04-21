package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Test 测试接口
// 用于测试服务器是否正常运行
func Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}
