package router

import (
	"backend/controllers"
	"backend/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有路由
// 返回配置好的Gin引擎实例
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 配置CORS，允许前端开发服务器访问
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175", "http://localhost:5176", "http://localhost:5177", "http://localhost:5178", "http://localhost:5179"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 认证相关路由（公开接口，无需登录）
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
		auth.POST("/generate-invitation-code", controllers.GenerateInvitationCode)
	}

	// API路由组
	api := r.Group("/api")
	{
		// ========== 公开接口（无需登录）==========

		// 文章相关（列表和详情无需认证，发布需要认证）
		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/:id", controllers.GetArticleByID)

		// 博客相关（列表、详情和搜索无需认证，发布、修改、删除需要认证）
		api.GET("/blogs", controllers.GetBlogs)
		api.GET("/blogs/search", controllers.SearchBlogs)
		api.GET("/blogs/:id", controllers.GetBlogByID)

		// 标签相关（只保留查询接口，标签的创建和删除由博客操作自动处理）
		api.GET("/tags", controllers.GetTags)
		api.GET("/tag-blogs/:tag", controllers.GetBlogsByTag)

		// 用户个人主页（公开访问）
		api.GET("/user/home/:homePath", controllers.GetUserByHomePath)

		// ========== 需要登录的接口 ==========
		// 应用认证中间件，以下接口需要有效的JWT Token
		api.Use(middlewares.AuthMiddleWare())
		{

			// 文章相关
			api.POST("/articles", controllers.CreateArticle)
			api.POST("/articles/:id/like", controllers.LikeArticle)
			api.GET("/articles/:id/like", controllers.GetArticleLikes)

			// 博客相关（增删改）
			api.POST("/blogs", controllers.CreateBlog)
			api.PUT("/blogs/:id", controllers.UpdateBlog)
			api.DELETE("/blogs/:id", controllers.DeleteBlog)
			api.POST("/blogs/:id/like", controllers.LikeBlog)

			// 用户相关
			api.GET("/user/profile", controllers.GetUserProfile)
			api.PUT("/user/profile", controllers.UpdateUserProfile)
			api.PUT("/user/password", controllers.UpdateUserPassword)

			// 管理员用户管理
			api.GET("/users", controllers.GetAllUsers)
			api.POST("/users", controllers.CreateUser)
			api.PUT("/users/:id", controllers.UpdateUser)
			api.DELETE("/users/:id", controllers.DeleteUser)

			// 邀请码管理
			api.GET("/invitation-codes", controllers.GetAllInvitationCodes)
			api.DELETE("/invitation-codes/:code", controllers.DeleteInvitationCode)
			api.DELETE("/invitation-codes/cleanup", controllers.CleanupExpiredCodesAPI)
		}
	}

	// 测试接口
	test := r.Group("/api/test")
	test.GET("/test", controllers.Test)

	return r
}
