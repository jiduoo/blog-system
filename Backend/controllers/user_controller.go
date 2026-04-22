package controllers

import (
	"backend/global"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserProfile 获取当前用户信息
// 需要登录后访问，返回用户的基本信息
func GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := global.Db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"homePath": user.HomePath,
		"isRoot":   user.IsRoot,
	})
}

// UpdateUserProfile 更新当前用户信息
// 需要登录后访问，可以更新个人主页路径
func UpdateUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		HomePath string `json:"homePath" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 检查HomePath是否已被其他用户使用
	var existingUser models.User
	if err := global.Db.Where("home_path = ? AND id != ?", input.HomePath, userID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Home path already exists"})
		return
	}

	// 更新用户个人主页路径
	if err := global.Db.Model(&models.User{}).Where("id = ?", userID).Update("home_path", input.HomePath).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// UpdateUserPassword 修改当前用户密码
// 需要登录后访问，需要提供旧密码和新密码
func UpdateUserPassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 获取用户信息
	var user models.User
	if err := global.Db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// 验证旧密码是否正确
	if !utils.CheckPassword(input.OldPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid old password"})
		return
	}

	// 使用BCrypt哈希新密码
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 更新密码
	if err := global.Db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// GetUserByHomePath 根据个人主页路径获取用户信息
// 公开接口，无需登录，用于访问用户的个人主页
func GetUserByHomePath(c *gin.Context) {
	homePath := c.Param("homePath")

	var user models.User
	if err := global.Db.Where("home_path = ?", homePath).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"homePath": user.HomePath,
	})
}

// GetAllUsers 获取所有用户列表
// 需要root用户权限
func GetAllUsers(c *gin.Context) {
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

	var users []models.User
	if err := global.Db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// CreateUser 创建新用户
// 需要root用户权限
func CreateUser(c *gin.Context) {
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

	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		IsRoot   bool   `json:"isRoot"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := global.Db.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 创建用户
	newUser := models.User{
		Username: input.Username,
		Password: hashedPassword,
		IsRoot:   input.IsRoot,
	}

	if err := global.Db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// UpdateUser 更新用户信息
// 需要root用户权限
func UpdateUser(c *gin.Context) {
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

	userID := c.Param("id")

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		IsRoot   *bool  `json:"isRoot"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var targetUser models.User
	if err := global.Db.First(&targetUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 更新用户信息
	updates := make(map[string]interface{})

	if input.Username != "" {
		// 检查新用户名是否已存在
		var existingUser models.User
		if err := global.Db.Where("username = ? AND id != ?", input.Username, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		updates["username"] = input.Username
	}

	if input.Password != "" {
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updates["password"] = hashedPassword
	}

	if input.IsRoot != nil {
		updates["is_root"] = *input.IsRoot
	}

	if len(updates) > 0 {
		if err := global.Db.Model(&targetUser).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser 删除用户
// 需要root用户权限
func DeleteUser(c *gin.Context) {
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

	userID := c.Param("id")

	// 不允许删除root用户
	var targetUser models.User
	if err := global.Db.First(&targetUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if targetUser.IsRoot {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete root user"})
		return
	}

	if err := global.Db.Delete(&targetUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
