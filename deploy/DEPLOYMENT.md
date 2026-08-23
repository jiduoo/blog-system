# 博客系统部署文档（Docker）

## 概述

本项目使用 Docker Compose 一键部署，包含以下服务：

| 服务 | 镜像 | 说明 |
| :--- | :--- | :--- |
| mysql | mysql:8.0 | 数据库 |
| redis | redis:7-alpine | 缓存 |
| backend | 自建（Go 1.22.5） | Gin + GORM 后端 API |
| frontend | 自建（Nginx） | Vue 静态资源 + API 反向代理 |

## 目录结构

```
项目根目录/
├── Backend/              # 后端源码
├── Frontend/             # 前端源码
└── deploy/
    ├── docker-compose.yml
    ├── Dockerfile.backend
    ├── Dockerfile.frontend
    ├── nginx.conf
    ├── config.yml         # Docker 环境后端配置（挂载进容器）
    └── DEPLOYMENT.md
```

> 注意：`docker-compose.yml` 的构建上下文为项目根目录（`..`），因此 Dockerfile 中使用 `Backend/`、`Frontend/` 相对路径。

## 部署步骤

### 1. 安装 Docker

```bash
curl -fsSL https://get.docker.com | bash -s docker
systemctl enable --now docker
```

### 2. 克隆项目

```bash
cd /opt
git clone <仓库地址> blog-system
cd blog-system
```

### 3. 启动服务

```bash
cd deploy
docker compose up -d --build
```

首次启动会构建镜像并初始化数据库，约需 1-2 分钟。

### 4. 验证部署

```bash
# 查看容器状态
docker compose ps

# 访问网站
curl http://localhost/api/blogs
```

浏览器打开 `http://<服务器IP>` 即可访问。

## 配置说明

### 数据库配置

默认账号密码定义在 [docker-compose.yml](docker-compose.yml) 的 `mysql` 服务环境变量中：

| 配置项 | 默认值 |
| :--- | :--- |
| 数据库名 | blog_system |
| 用户名 | blog_user |
| 密码 | blog_password |
| root 密码 | root_password |

如需修改，请同步修改 [config.yml](config.yml) 中的 `dsn` 字段。

### 后端配置

[config.yml](config.yml) 以只读卷形式挂载到后端容器 `/app/config/config.yml`，修改后重启 backend 容器即可生效：

```bash
docker compose restart backend
```

### 端口

仅对外暴露 **80** 端口（前端 Nginx）。MySQL、Redis、后端 API 仅在 Docker 内网通信，提升安全性。

## 常用命令

```bash
# 启动
docker compose up -d --build

# 停止
docker compose down

# 重启
docker compose restart

# 查看日志
docker compose logs -f
docker compose logs -f backend

# 查看状态
docker compose ps

# 重新构建（代码更新后）
docker compose up -d --build backend
```

## 默认账号

| 类型 | 用户名 | 密码 |
| :--- | :--- | :--- |
| 管理员 | root | root |

**首次登录后请立即修改密码。**

## 数据持久化

数据存储在 Docker 卷中：

| 卷 | 用途 |
| :--- | :--- |
| deploy_mysql_data | MySQL 数据 |
| deploy_redis_data | Redis 数据 |
| deploy_backend_logs | 后端日志 |

备份数据库：

```bash
docker exec blog-mysql mysqldump -u root -proot_password blog_system > backup.sql
```

## 常见问题

### 1. 端口 80 被占用

```bash
# 查看占用
netstat -tlnp | grep :80
# 修改 docker-compose.yml 中 frontend 的端口映射
```

### 2. 后端启动失败

```bash
docker compose logs backend
```

常见原因：MySQL 未就绪（等待 healthcheck 通过）、config.yml 配置错误。

### 3. 重新初始化数据库

```bash
docker compose down -v   # 注意：-v 会删除数据卷
docker compose up -d --build
```
