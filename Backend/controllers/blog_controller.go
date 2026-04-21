package controllers

import (
	"backend/global"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateBlog 创建新博客
// 需要登录后访问，创建博客并关联标签
func CreateBlog(c *gin.Context) {
	var input struct {
		Title   string   `json:"title" binding:"required"`
		Content string   `json:"content" binding:"required"`
		Preview string   `json:"preview" binding:"required"`
		Tags    []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从认证中间件获取当前用户名
	username, _ := c.Get("username")

	blog := models.Blog{
		Title:   input.Title,
		Content: input.Content,
		Preview: input.Preview,
		Author:  username.(string),
	}

	// 自动迁移博客表
	if err := global.Db.AutoMigrate(&blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建博客记录
	if err := global.Db.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 处理标签关联
	if len(input.Tags) > 0 {
		for _, tagName := range input.Tags {
			var tag models.Tag
			// 查找或创建标签
			if err := global.Db.Where("name = ?", tagName).FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle tags: " + err.Error()})
				return
			}
			// 关联标签到博客
			if err := global.Db.Model(&blog).Association("Tags").Append(&tag); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate tags: " + err.Error()})
				return
			}
		}
	}

	// 重新加载博客以包含标签
	global.Db.Preload("Tags").First(&blog, blog.ID)

	c.JSON(http.StatusCreated, blog)
}

// GetBlogs 获取所有博客列表
// 公开接口，无需登录，返回所有博客及其标签
func GetBlogs(c *gin.Context) {
	var blogs []models.Blog
	result := global.Db.Preload("Tags").Find(&blogs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, blogs)
}

// GetBlogByID 获取指定ID的博客详情
// 公开接口，无需登录，同时增加浏览量
func GetBlogByID(c *gin.Context) {
	id := c.Param("id")
	var blog models.Blog
	result := global.Db.Preload("Tags").First(&blog, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// 增加浏览量
	blog.Views++
	global.Db.Save(&blog)

	c.JSON(http.StatusOK, blog)
}

// LikeBlog 点赞博客
// 公开接口，无需登录，每次调用增加1个点赞
func LikeBlog(c *gin.Context) {
	id := c.Param("id")
	var blog models.Blog
	if err := global.Db.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// 增加点赞数
	blog.Likes++
	if err := global.Db.Save(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully liked the blog", "likes": blog.Likes})
}

// UpdateBlog 更新博客
// 需要登录后访问，只能修改自己的博客
func UpdateBlog(c *gin.Context) {
	id := c.Param("id")
	username, _ := c.Get("username")

	var blog models.Blog
	if err := global.Db.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// 检查是否是博客作者
	if blog.Author != username.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own blogs"})
		return
	}

	var input struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Preview string   `json:"preview"`
		Tags    []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	if input.Title != "" {
		blog.Title = input.Title
	}
	if input.Content != "" {
		blog.Content = input.Content
	}
	if input.Preview != "" {
		blog.Preview = input.Preview
	}

	if err := global.Db.Save(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新标签关联
	if input.Tags != nil {
		// 清除旧标签关联
		global.Db.Model(&blog).Association("Tags").Clear()

		// 添加新标签
		for _, tagName := range input.Tags {
			var tag models.Tag
			if err := global.Db.Where("name = ?", tagName).FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
				continue
			}
			global.Db.Model(&blog).Association("Tags").Append(&tag)
		}
	}

	// 重新加载博客以包含标签
	global.Db.Preload("Tags").First(&blog, blog.ID)

	c.JSON(http.StatusOK, blog)
}

// DeleteBlog 删除博客
// 需要登录后访问，只能删除自己的博客
func DeleteBlog(c *gin.Context) {
	id := c.Param("id")
	username, _ := c.Get("username")

	var blog models.Blog
	if err := global.Db.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// 检查是否是博客作者
	if blog.Author != username.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own blogs"})
		return
	}

	// 清除标签关联
	global.Db.Model(&blog).Association("Tags").Clear()

	// 删除博客
	if err := global.Db.Delete(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}

// SearchBlogs 搜索博客
// 公开接口，无需登录，支持按标题和内容搜索
func SearchBlogs(c *gin.Context) {
	keyword := c.Query("keyword")

	var blogs []models.Blog
	query := global.Db.Preload("Tags")

	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ? OR author LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if err := query.Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blogs)
}
