package controllers

import (
	"backend/global"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTags 获取所有有博客关联的标签列表
// 公开接口，无需登录，只返回至少有一篇博客使用的标签
func GetTags(c *gin.Context) {
	var tags []models.Tag
	// 只查询至少有一篇博客使用的标签
	result := global.Db.Joins("JOIN blog_tags ON tags.id = blog_tags.tag_id").Group("tags.id").Find(&tags)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, tags)
}

// GetTagByID 根据ID获取标签详情
// 公开接口，无需登录
func GetTagByID(c *gin.Context) {
	id := c.Param("id")

	var tag models.Tag
	if err := global.Db.First(&tag, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// GetBlogsByTag 根据标签名获取关联的博客列表
// 公开接口，无需登录
func GetBlogsByTag(c *gin.Context) {
	tagName := c.Param("tag")

	// 查找标签
	var tag models.Tag
	if err := global.Db.Where("name = ?", tagName).First(&tag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// 查找该标签关联的所有博客，并预加载标签信息
	var blogs []models.Blog
	if err := global.Db.Model(&tag).Association("Blogs").Find(&blogs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 为每个博客预加载标签
	for i := range blogs {
		global.Db.Preload("Tags").First(&blogs[i], blogs[i].ID)
	}

	c.JSON(http.StatusOK, blogs)
}
