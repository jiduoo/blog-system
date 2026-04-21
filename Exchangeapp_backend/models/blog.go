package models

import (
	"gorm.io/gorm"
)

// Blog 博客模型
// 用于存储博客文章信息
type Blog struct {
	gorm.Model        // 包含ID、CreatedAt、UpdatedAt、DeletedAt字段
	Title    string   `json:"title" binding:"required"` // 博客标题
	Content  string   `json:"content" binding:"required"` // 博客正文内容
	Preview  string   `json:"preview" binding:"required"` // 预览文本，用于列表展示
	Author   string   `json:"author"` // 作者用户名
	Views    int      `json:"views" gorm:"default:0"` // 浏览量
	Likes    int      `json:"likes" gorm:"default:0"` // 点赞数
	Tags     []Tag    `json:"tags" gorm:"many2many:blog_tags;"` // 关联的标签，多对多关系
}

// TableName 指定表名为 blogs
func (Blog) TableName() string {
	return "blogs"
}
