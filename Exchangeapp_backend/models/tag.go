package models

import (
	"gorm.io/gorm"
)

// Tag 标签模型
// 用于对博客文章进行分类和标记
type Tag struct {
	gorm.Model        // 包含ID、CreatedAt、UpdatedAt、DeletedAt字段
	Name  string `gorm:"unique;not null"` // 标签名称，唯一且不能为空
	Blogs []Blog `gorm:"many2many:blog_tags;"` // 关联的博客文章，多对多关系
}

// TableName 指定表名为 tags
func (Tag) TableName() string {
	return "tags"
}
