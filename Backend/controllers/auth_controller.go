package controllers

import (
	"backend/global"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
// 公开接口，无需登录，需要提供用户名、密码和有效的注册码
// 注册成功后返回JWT Token
func Register(ctx *gin.Context) {
	var input struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		InvitationCode string `json:"invitationCode"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查注册码是否为空
	if input.InvitationCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "注册码不能为空"})
		return
	}

	// 验证注册码是否有效
	valid, err := ValidateInvitationCode(input.InvitationCode)
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用BCrypt哈希密码
	hashedPwd, err := utils.HashPassword(input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建用户对象
	user := models.User{
		Username: input.Username,
		Password: hashedPwd,
		IsRoot:   false,
	}

	// 生成JWT Token
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 自动迁移用户表
	if err := global.Db.AutoMigrate(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建用户记录
	if err := global.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 标记注册码已使用
	if err := MarkInvitationCodeAsUsed(input.InvitationCode, user.Username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "注册成功但注册码标记失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// Login 用户登录
// 公开接口，无需登录，需要提供用户名和密码
// 登录成功后返回JWT Token
func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
		return
	}

	// 验证密码
	if !utils.CheckPassword(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
		return
	}

	// 生成JWT Token
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
