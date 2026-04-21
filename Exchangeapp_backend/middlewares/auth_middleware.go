package middlewares

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleWare 验证用户身份
// 1. 从请求头获取 Authorization Token
// 2. 解析 Token 获取用户名
// 3. 根据用户名查询用户信息，设置用户ID到context
func AuthMiddleWare() gin.HandlerFunc{
	return func(ctx *gin.Context){
		// 获取请求头中的 Authorization Token
		token := ctx.GetHeader("Authorization")
		if token == ""{
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			ctx.Abort()
			return
		}

		// 解析 Token 获取用户名
		username, err := utils.ParseJWT(token)
		if err !=nil{
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// 根据用户名查询用户信息
		var user models.User
		if err := global.Db.Where("username = ?", username).First(&user).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			ctx.Abort()
			return
		}

		// 设置用户名和用户ID到context，供后续处理器使用
		ctx.Set("username", username)
		ctx.Set("userID", user.ID)
		ctx.Next()
	}
}
