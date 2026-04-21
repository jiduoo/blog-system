# 博客系统

一个基于 Vue 3 + TypeScript + Go + Gin + GORM 的前后端分离博客系统，支持无需登录查看文章，登录后可发布、管理文章，包含标签分类、时间归档、个人主页、代码高亮、Markdown 编辑器等功能。

## 📋 功能特性

- **无需登录查看**：访客可直接浏览所有文章
- **用户系统**：支持注册、登录、个人中心
- **文章管理**：发布、编辑、删除文章
- **标签系统**：文章标签分类，支持按标签筛选
- **时间归档**：按年份和日期组织博客
- **个人主页**：每个用户可设置唯一的个人主页路径
- **Markdown 编辑器**：支持 Markdown 写作，实时预览
- **代码高亮**：代码块支持语法高亮，带复制、下载、截图功能
- **响应式设计**：适配不同设备屏幕

## 🛠 技术栈

### 前端
- Vue 3 + TypeScript
- Vue Router + Pinia (状态管理)
- Axios (HTTP请求)
- Element Plus (UI组件)
- Markdown-it (Markdown渲染)
- Highlight.js (代码高亮)
- HTML2Canvas (截图功能)
- Vite (构建工具)

### 后端
- Go 1.22.5
- Gin 框架
- GORM (ORM)
- MySQL (数据库)
- JWT (认证)
- Viper (配置管理)

## 📁 项目结构

```
├── Backend/                  # 后端 Go 代码
│   ├── config/               # 配置文件
│   ├── controllers/          # 控制器
│   ├── middlewares/          # 中间件
│   ├── models/               # 数据模型
│   ├── router/               # 路由配置
│   ├── utils/                # 工具函数
│   └── main.go               # 主入口
├── Frontend/                 # 前端 Vue 代码
│   ├── src/                  # 源代码
│   │   ├── components/       # 组件
│   │   ├── router/           # 路由
│   │   ├── store/            # 状态管理
│   │   ├── views/            # 页面
│   │   ├── App.vue           # 根组件
│   │   └── main.ts           # 主入口
│   └── package.json          # 依赖配置
├── .gitignore                # Git 忽略文件
├── LICENSE                   # 许可证
└── README.md                 # 项目说明
```

## 🔧 快速开始

### 前置要求
- Go 1.22.5+
- Node.js 18+
- MySQL 5.7+
- Redis (可选，用于缓存)

### 后端启动

1. **配置数据库**
   - 修改 `Backend/config/config.yml` 中的数据库连接信息

2. **安装依赖**
   ```bash
   cd Backend
   go mod download
   ```

3. **启动服务**
   ```bash
   go run main.go
   ```
   服务默认运行在 `http://localhost:3000`

### 前端启动

1. **安装依赖**
   ```bash
   cd Frontend
   npm install
   ```

2. **启动开发服务器**
   ```bash
   npm run dev
   ```
   前端默认运行在 `http://localhost:5173`

## 🗄 数据库配置

### 配置文件
修改 `Backend/config/config.yml`：

```yaml
database:
  dsn: root:123456@tcp(192.168.30.10:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
  MaxIdleConns: 11
  MaxOpenCons: 114

redis:
  addr: 192.168.30.10:6379
  DB: 0
  Password: ""
```

### 数据库迁移
项目启动时会自动创建表结构，无需手动迁移。

## 🌐 主要功能

### 1. 博客列表
- 支持按时间排序
- 支持按标签筛选
- 支持时间归档
- 支持搜索功能

### 2. 博客详情
- Markdown 渲染
- 代码高亮
- 代码块操作（复制、下载、截图）
- 目录导航
- 点赞功能

### 3. 写博客
- Markdown 编辑器
- 实时预览
- 标签管理
- 下载 Markdown 文件

### 4. 个人中心
- 个人信息管理
- 密码修改
- 个人主页路径设置

### 5. 后台管理
- 文章管理
- 标签管理
- 用户管理

## 📝 API 文档

### 公共接口
- `GET /blogs` - 获取博客列表
- `GET /blogs/:id` - 获取博客详情
- `GET /tags` - 获取所有标签
- `GET /blogs/tag/:tag` - 按标签获取博客

### 需要认证的接口
- `POST /blogs` - 创建博客
- `PUT /blogs/:id` - 更新博客
- `DELETE /blogs/:id` - 删除博客
- `POST /blogs/:id/like` - 点赞博客
- `PUT /users/profile` - 更新用户信息
- `PUT /users/password` - 修改密码

## 🚀 部署

### 前端构建
```bash
cd Frontend
npm run build
```
构建产物位于 `dist` 目录，可部署到任何静态文件服务器。

### 后端部署
1. **编译**
   ```bash
   cd Backend
   go build -o blog-backend .
   ```

2. **运行**
   ```bash
   ./blog-backend
   ```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🌟 致谢

- [Vue.js](https://vuejs.org/)
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [Markdown-it](https://markdown-it.github.io/)
- [Highlight.js](https://highlightjs.org/)

---

**作者**: jiduoo
**GitHub**: [https://github.com/jiduoo/blog-system](https://github.com/jiduoo/blog-system)