package models

import (
	"time"

	"gorm.io/gorm"
)

// InvitationCode 注册码模型
// 用于用户注册时的邀请码验证
type InvitationCode struct {
	gorm.Model        // 包含ID、CreatedAt、UpdatedAt、DeletedAt字段
	Code      string    `gorm:"unique;not null"` // 注册码，唯一且不能为空
	Used      bool      `gorm:"default:false"` // 是否已被使用
	UsedBy    string    `gorm:"size:255"` // 使用者用户名
	CreatedBy string    `gorm:"size:255;not null"` // 创建者用户名
	ExpiresAt time.Time `gorm:"not null"` // 过期时间
}

// TableName 指定表名为 invitation_codes
func (InvitationCode) TableName() string {
	return "invitation_codes"
}
