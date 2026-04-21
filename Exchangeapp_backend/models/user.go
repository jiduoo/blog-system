package models

import "gorm.io/gorm"

// User 用户模型
// 用于存储用户账户信息
type User struct {
	gorm.Model        // 包含ID、CreatedAt、UpdatedAt、DeletedAt字段
	Username string   `gorm:"unique"` // 用户名，唯一索引
	Password string   // 密码，BCrypt哈希存储
	IsRoot   bool     `gorm:"default:false"` // 是否为root管理员
	HomePath string   `gorm:"unique"` // 个人主页路径，唯一索引，如 /zhangtest
}
