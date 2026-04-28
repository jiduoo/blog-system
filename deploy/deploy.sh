#!/bin/bash

# ============================================
# 博客系统部署脚本
# 适用环境：CentOS 7/8 / Ubuntu 18.04+
# 使用方式：bash deploy.sh [start|stop|restart|status|install]
# ============================================

# 配置参数
APP_NAME="blog-system"
APP_DIR="/opt/blog-system"
BACKEND_PORT=3000
FRONTEND_PORT=80
MYSQL_PORT=3306
REDIS_PORT=6379

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 安装依赖
install_dependencies() {
    log_info "开始安装依赖..."
    
    # 更新系统
    if [ -f /etc/redhat-release ]; then
        # CentOS/RHEL
        log_info "检测到 CentOS/RHEL 系统"
        yum update -y
        yum install -y wget curl git unzip
    elif [ -f /etc/debian_version ]; then
        # Ubuntu/Debian
        log_info "检测到 Ubuntu/Debian 系统"
        apt-get update -y
        apt-get install -y wget curl git unzip
    else
        log_error "不支持的操作系统"
        exit 1
    fi

    log_info "依赖安装完成"
}

# 安装 Node.js
install_nodejs() {
    log_info "开始安装 Node.js..."
    
    if command -v node &> /dev/null; then
        log_warn "Node.js 已安装，跳过安装步骤"
        return
    fi

    # 安装 Node.js 18.x
    curl -fsSL https://rpm.nodesource.com/setup_18.x | bash -
    if [ -f /etc/redhat-release ]; then
        yum install -y nodejs
    else
        apt-get install -y nodejs
    fi

    log_info "Node.js 安装完成"
    node --version
    npm --version
}

# 安装 Go
install_go() {
    log_info "开始安装 Go..."
    
    if command -v go &> /dev/null; then
        log_warn "Go 已安装，跳过安装步骤"
        return
    fi

    # 安装 Go 1.22.5
    wget https://dl.google.com/go/go1.22.5.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    source /etc/profile

    log_info "Go 安装完成"
    go version
}

# 安装 MySQL
install_mysql() {
    log_info "开始安装 MySQL..."
    
    if command -v mysql &> /dev/null; then
        log_warn "MySQL 已安装，跳过安装步骤"
        return
    fi

    if [ -f /etc/redhat-release ]; then
        # CentOS/RHEL
        yum install -y mariadb-server mariadb
        systemctl enable mariadb
        systemctl start mariadb
    else
        # Ubuntu/Debian
        apt-get install -y mysql-server
        systemctl enable mysql
        systemctl start mysql
    fi

    log_info "MySQL 安装完成"
}

# 安装 Redis
install_redis() {
    log_info "开始安装 Redis..."
    
    if command -v redis-server &> /dev/null; then
        log_warn "Redis 已安装，跳过安装步骤"
        return
    fi

    if [ -f /etc/redhat-release ]; then
        yum install -y redis
    else
        apt-get install -y redis-server
    fi
    
    systemctl enable redis
    systemctl start redis

    log_info "Redis 安装完成"
}

# 配置数据库
configure_database() {
    log_info "开始配置数据库..."
    
    # 创建数据库
    mysql -u root <<EOF
CREATE DATABASE IF NOT EXISTS blog_system DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'blog_user'@'localhost' IDENTIFIED BY 'blog_password';
GRANT ALL PRIVILEGES ON blog_system.* TO 'blog_user'@'localhost';
FLUSH PRIVILEGES;
EOF

    if [ $? -eq 0 ]; then
        log_info "数据库配置完成"
    else
        log_error "数据库配置失败"
        exit 1
    fi
}

# 部署后端
deploy_backend() {
    log_info "开始部署后端..."
    
    # 获取当前脚本所在目录
    SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
    PROJECT_DIR=$(dirname "$SCRIPT_DIR")
    
    # 创建应用目录
    mkdir -p $APP_DIR/Backend
    
    # 复制后端代码（只复制源文件，跳过目标目录）
    if [ "$PROJECT_DIR/Backend" != "$APP_DIR/Backend" ]; then
        rsync -av --delete "$PROJECT_DIR/Backend/" "$APP_DIR/Backend/" --exclude='*.log' --exclude='__pycache__'
    else
        log_warn "后端代码已在目标位置，跳过复制"
    fi
    
    # 修改配置文件
    sed -i "s|root:123456@tcp(192.168.30.10:3306)/test|blog_user:blog_password@tcp(localhost:3306)/blog_system|g" $APP_DIR/Backend/config/config.yml
    
    # 编译后端
    cd $APP_DIR/Backend
    go mod download
    go build -o blog-backend .
    
    if [ $? -eq 0 ]; then
        log_info "后端部署完成"
    else
        log_error "后端编译失败"
        exit 1
    fi
}

# 部署前端
deploy_frontend() {
    log_info "开始部署前端..."
    
    # 获取当前脚本所在目录
    SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
    PROJECT_DIR=$(dirname "$SCRIPT_DIR")
    
    # 创建应用目录
    mkdir -p $APP_DIR/Frontend
    
    # 复制前端代码（只复制源文件，跳过目标目录）
    if [ "$PROJECT_DIR/Frontend" != "$APP_DIR/Frontend" ]; then
        rsync -av --delete "$PROJECT_DIR/Frontend/" "$APP_DIR/Frontend/" --exclude='node_modules' --exclude='dist' --exclude='*.log'
    else
        log_warn "前端代码已在目标位置，跳过复制"
    fi
    
    # 安装依赖并构建
    cd $APP_DIR/Frontend
    
    # 使用国内npm镜像
    npm config set registry https://registry.npmmirror.com/
    
    # 安装依赖
    npm install
    
    # 设置vue-tsc执行权限
    chmod +x node_modules/.bin/vue-tsc 2>/dev/null || true
    
    # 构建
    npm run build
    
    if [ $? -eq 0 ]; then
        log_info "前端部署完成"
    else
        log_error "前端构建失败"
        exit 1
    fi
}

# 配置 Nginx
configure_nginx() {
    log_info "开始配置 Nginx..."
    
    if ! command -v nginx &> /dev/null; then
        if [ -f /etc/redhat-release ]; then
            yum install -y nginx
        else
            apt-get install -y nginx
        fi
    fi

    # 创建 Nginx 配置
    cat > /etc/nginx/conf.d/blog.conf <<EOF
server {
    listen 80;
    server_name localhost;

    # 前端静态文件
    location / {
        root $APP_DIR/Frontend/dist;
        index index.html;
        try_files \$uri \$uri/ /index.html;
    }

    # 后端 API 代理
    location /api/ {
        proxy_pass http://localhost:$BACKEND_PORT/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }
}
EOF

    # 重启 Nginx
    systemctl enable nginx
    systemctl restart nginx

    log_info "Nginx 配置完成"
}

# 创建系统服务
create_service() {
    log_info "创建系统服务..."
    
    cat > /etc/systemd/system/blog-backend.service <<EOF
[Unit]
Description=Blog System Backend
After=network.target mysqld.service redis.service

[Service]
Type=simple
User=root
WorkingDirectory=$APP_DIR/Backend
ExecStart=$APP_DIR/Backend/blog-backend
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable blog-backend

    log_info "系统服务创建完成"
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    systemctl start blog-backend
    
    if [ $? -eq 0 ]; then
        log_info "服务启动成功"
    else
        log_error "服务启动失败"
        exit 1
    fi
}

# 停止服务
stop_services() {
    log_info "停止服务..."
    
    systemctl stop blog-backend
    
    if [ $? -eq 0 ]; then
        log_info "服务停止成功"
    else
        log_error "服务停止失败"
        exit 1
    fi
}

# 重启服务
restart_services() {
    log_info "重启服务..."
    
    systemctl restart blog-backend
    
    if [ $? -eq 0 ]; then
        log_info "服务重启成功"
    else
        log_error "服务重启失败"
        exit 1
    fi
}

# 查看状态
check_status() {
    log_info "查看服务状态..."
    
    systemctl status blog-backend
    
    echo ""
    log_info "检查端口占用情况:"
    netstat -tlnp | grep -E ":$BACKEND_PORT|:$FRONTEND_PORT"
}

# 安装模式
install_mode() {
    log_info "执行安装模式..."
    
    install_dependencies
    install_nodejs
    install_go
    install_mysql
    install_redis
    configure_database
    deploy_backend
    deploy_frontend
    configure_nginx
    create_service
    start_services
    
    log_info "安装完成！"
    log_info "请访问 http://<服务器IP> 查看博客系统"
    log_info "默认管理员账号: root / root"
}

# 主函数
case "$1" in
    install)
        install_mode
        ;;
    start)
        start_services
        ;;
    stop)
        stop_services
        ;;
    restart)
        restart_services
        ;;
    status)
        check_status
        ;;
    *)
        echo "使用方式: $0 {install|start|stop|restart|status}"
        echo ""
        echo "参数说明:"
        echo "  install   - 完整安装并部署项目"
        echo "  start     - 启动服务"
        echo "  stop      - 停止服务"
        echo "  restart   - 重启服务"
        echo "  status    - 查看服务状态"
        exit 1
        ;;
esac

exit 0