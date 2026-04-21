package models

import "gorm.io/gorm"

// Article 文章模型
// 用于存储新闻或文章信息（与博客不同，这是另一个模块）
type Article struct {
	gorm.Model        // 包含ID、CreatedAt、UpdatedAt、DeletedAt字段
	Title string `binding:"required"` // 文章标题
	Content string `binding:"required"` // 文章正文内容
	Preview string `binding:"required"` // 预览文本
}
