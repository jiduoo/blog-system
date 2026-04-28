# 博客系统部署文档

## 概述

本文档详细介绍了如何在阿里云 ECS 服务器上部署博客系统。提供两种部署方式：

1. **传统部署方式**：使用 Shell 脚本直接在服务器上部署
2. **容器化部署方式**：使用 Docker Compose 进行容器化部署

## 前置准备

### 服务器要求

| 配置项 | 最低配置 | 推荐配置 |
| :--- | :--- | :--- |
| CPU | 1核 | 2核及以上 |
| 内存 | 2GB | 4GB及以上 |
| 存储 | 40GB | 100GB及以上 |
| 操作系统 | CentOS 7/8 或 Ubuntu 18.04+ | CentOS 8 |

### 网络配置

确保阿里云 ECS 安全组开放以下端口：

- **80**：HTTP 服务端口
- **3000**：后端 API 端口（可选，用于直接访问）
- **3306**：MySQL 端口（可选，用于远程连接）
- **6379**：Redis 端口（可选，用于远程连接）

### 域名配置（可选）

如果需要使用域名访问，提前在阿里云域名解析中添加 A 记录指向服务器 IP。

---

## 方式一：传统部署

### 1. 登录服务器

```bash
ssh root@<服务器IP>
```

### 2. 安装 Git

```bash
# CentOS/RHEL
yum install -y git

# Ubuntu/Debian
apt-get install -y git
```

### 3. 克隆项目

```bash
cd /opt
git clone https://github.com/jiduoo/blog-system.git
cd blog-system
```

### 4. 运行部署脚本

```bash
cd deploy
chmod +x deploy.sh
bash deploy.sh install
```

### 5. 脚本参数说明

```bash
# 安装并部署（首次部署使用）
bash deploy.sh install

# 启动服务
bash deploy.sh start

# 停止服务
bash deploy.sh stop

# 重启服务
bash deploy.sh restart

# 查看服务状态
bash deploy.sh status
```

### 6. 手动配置（可选）

如果脚本执行失败，可以手动执行以下步骤：

#### 6.1 安装依赖

```bash
# CentOS/RHEL
yum install -y wget curl git unzip nginx

# Ubuntu/Debian
apt-get install -y wget curl git unzip nginx
```

#### 6.2 安装 Node.js 18.x

```bash
curl -fsSL https://rpm.nodesource.com/setup_18.x | bash -
yum install -y nodejs
```

#### 6.3 安装 Go 1.22.5

```bash
wget https://dl.google.com/go/go1.22.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile
```

#### 6.4 配置数据库

```bash
mysql -u root
```

```sql
CREATE DATABASE blog_system DEFAULT CHARACTER SET utf8mb4;
CREATE USER 'blog_user'@'localhost' IDENTIFIED BY 'blog_password';
GRANT ALL PRIVILEGES ON blog_system.* TO 'blog_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

#### 6.5 编译后端

```bash
cd /opt/blog-system/Backend
go mod download
go build -o blog-backend .
```

#### 6.6 构建前端

```bash
cd /opt/blog-system/Frontend
npm install
npm run build
```

#### 6.7 配置 Nginx

创建 `/etc/nginx/conf.d/blog.conf`：

```nginx
server {
    listen 80;
    server_name localhost;

    location / {
        root /opt/blog-system/Frontend/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://localhost:3000/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

```bash
systemctl restart nginx
```

#### 6.8 创建系统服务

创建 `/etc/systemd/system/blog-backend.service`：

```ini
[Unit]
Description=Blog System Backend
After=network.target mysqld.service redis.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/blog-system/Backend
ExecStart=/opt/blog-system/Backend/blog-backend
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable blog-backend
systemctl start blog-backend
```

---

## 方式二：容器化部署

### 1. 安装 Docker 和 Docker Compose

```bash
# 安装 Docker
curl -fsSL https://get.docker.com | bash -s docker

# 安装 Docker Compose
curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# 启动 Docker 服务
systemctl enable docker
systemctl start docker
```

### 2. 克隆项目

```bash
cd /opt
git clone https://github.com/jiduoo/blog-system.git
cd blog-system
```

### 3. 配置环境变量（可选）

编辑 `deploy/docker-compose.yml`，根据需要修改以下配置：

- MySQL 数据库密码
- Redis 配置
- 端口映射

### 4. 启动容器

```bash
cd deploy
docker-compose up -d
```

### 5. 常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 查看日志
docker-compose logs -f

# 查看运行状态
docker-compose ps

# 重启服务
docker-compose restart

# 构建镜像（修改代码后）
docker-compose build
```

### 6. 数据持久化

容器化部署时，数据会自动持久化到以下目录：

- MySQL 数据：`/var/lib/docker/volumes/deploy_mysql_data/_data`
- Redis 数据：`/var/lib/docker/volumes/deploy_redis_data/_data`
- 后端日志：`/var/lib/docker/volumes/deploy_backend_logs/_data`

---

## 验证部署

### 1. 检查服务状态

```bash
# 传统部署
bash deploy.sh status

# 容器化部署
docker-compose ps
```

### 2. 访问网站

打开浏览器访问：

```
http://<服务器IP>
```

### 3. 测试 API

```bash
curl http://<服务器IP>/api/blogs
```

---

## 默认账号

| 账号类型 | 用户名 | 密码 |
| :--- | :--- | :--- |
| 管理员 | root | root |

**注意**：首次登录后请立即修改密码！

---

## 常见问题

### 1. 端口被占用

```bash
# 查看端口占用
netstat -tlnp | grep 80
netstat -tlnp | grep 3000

# 杀死占用进程
kill -9 <PID>
```

### 2. Nginx 配置错误

```bash
# 检查配置文件
nginx -t

# 查看错误日志
cat /var/log/nginx/error.log
```

### 3. 后端服务启动失败

```bash
# 查看日志（传统部署）
journalctl -u blog-backend -f

# 查看日志（容器化部署）
docker-compose logs -f backend
```

### 4. 数据库连接失败

确保 MySQL 服务已启动：

```bash
# 传统部署
systemctl status mariadb  # CentOS
systemctl status mysql    # Ubuntu

# 容器化部署
docker-compose logs mysql
```

---

## 安全建议

1. **修改默认密码**：登录后立即修改 root 用户密码
2. **关闭不必要的端口**：在阿里云安全组中只开放必要端口
3. **启用 HTTPS**：配置 SSL 证书（推荐使用 Let's Encrypt）
4. **定期备份**：定期备份数据库和重要文件
5. **更新系统**：定期更新系统和依赖包

---

## 维护命令

```bash
# 查看系统状态
top
htop

# 查看磁盘使用
df -h

# 查看内存使用
free -h

# 查看网络连接
netstat -tlnp
ss -tlnp
```

---

## 联系信息

如有问题，请联系：

- 作者：jiduoo
- GitHub：https://github.com/jiduoo/blog-system