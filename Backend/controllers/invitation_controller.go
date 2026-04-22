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
	// 从context中获取用户名
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 检查用户是否存在
	var user models.User
	if err := global.Db.Where("username = ?", username).First(&user).Error; err != nil {
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
		CreatedBy: username.(string),
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

// MarkInvitationCodeAsUsed 标记注册码为已使用并删除
// 注册成功后调用，记录使用者和使用时间，然后删除注册码
func MarkInvitationCodeAsUsed(code string, username string) error {
	return global.Db.Where("code = ?", code).Delete(&models.InvitationCode{}).Error
}

// generateRandomCode 生成指定长度的随机十六进制字符串
func generateRandomCode(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}

// GetAllInvitationCodes 获取所有邀请码
// 需要root用户权限
func GetAllInvitationCodes(c *gin.Context) {
	// 检查是否是root用户
	username, _ := c.Get("username")
	var user models.User
	if err := global.Db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !user.IsRoot {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only root user can access this endpoint"})
		return
	}

	var invitationCodes []models.InvitationCode
	if err := global.Db.Find(&invitationCodes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get invitation codes"})
		return
	}

	c.JSON(http.StatusOK, invitationCodes)
}

// DeleteInvitationCode 删除邀请码
// 需要root用户权限
func DeleteInvitationCode(c *gin.Context) {
	// 检查是否是root用户
	username, _ := c.Get("username")
	var user models.User
	if err := global.Db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !user.IsRoot {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only root user can access this endpoint"})
		return
	}

	code := c.Param("code")

	if err := global.Db.Where("code = ?", code).Delete(&models.InvitationCode{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invitation code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation code deleted successfully"})
}

// CleanupExpiredCodes 清理过期的邀请码
// 可以定时调用或在获取邀请码列表时自动清理
func CleanupExpiredCodes() {
	global.Db.Where("expires_at < ?", time.Now()).Delete(&models.InvitationCode{})
	global.Db.Where("used = ?", true).Delete(&models.InvitationCode{})
}

// CleanupExpiredCodesAPI 清理过期的邀请码API
// 需要root用户权限
func CleanupExpiredCodesAPI(c *gin.Context) {
	// 检查是否是root用户
	username, _ := c.Get("username")
	var user models.User
	if err := global.Db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !user.IsRoot {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only root user can access this endpoint"})
		return
	}

	CleanupExpiredCodes()
	c.JSON(http.StatusOK, gin.H{"message": "Expired invitation codes cleaned up successfully"})
}
