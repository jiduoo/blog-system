package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"backend/global"
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateInvitationCode 生成新的注册码
// 需要root用户权限，注册码有效期为7天
func GenerateInvitationCode(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否存在
	var user models.User
	if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	// 只有root用户可以创建注册码
	if !user.IsRoot {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只有root用户可以创建注册码"})
		return
	}

	// 生成16位随机注册码
	code := generateRandomCode(16)

	invitationCode := models.InvitationCode{
		Code:      code,
		Used:      false,
		CreatedBy: input.Username,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7天有效期
	}

	// 自动迁移注册码表
	if err := global.Db.AutoMigrate(&invitationCode); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建注册码记录
	if err := global.Db.Create(&invitationCode).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": code})
}

// ValidateInvitationCode 验证注册码是否有效
// 检查注册码是否存在、是否已使用、是否过期
func ValidateInvitationCode(code string) (bool, error) {
	var invitationCode models.InvitationCode

	if err := global.Db.Where("code = ?", code).First(&invitationCode).Error; err != nil {
		return false, errors.New("注册码不存在")
	}

	if invitationCode.Used {
		return false, errors.New("注册码已被使用")
	}

	if time.Now().After(invitationCode.ExpiresAt) {
		return false, errors.New("注册码已过期")
	}

	return true, nil
}

// MarkInvitationCodeAsUsed 标记注册码为已使用
// 注册成功后调用，记录使用者和使用时间
func MarkInvitationCodeAsUsed(code string, username string) error {
	return global.Db.Model(&models.InvitationCode{}).Where("code = ?", code).Updates(map[string]interface{}{
		"used":   true,
		"usedBy": username,
	}).Error
}

// generateRandomCode 生成指定长度的随机十六进制字符串
func generateRandomCode(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}
