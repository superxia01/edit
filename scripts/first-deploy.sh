#!/bin/bash
# ============================================
# Edit Business 首次部署脚本
# ============================================
#
# 使用说明：
# 1. 确保 scripts/ops/config.sh 已配置
# 2. 添加可执行权限：chmod +x scripts/first-deploy.sh
# 3. 执行首次部署：./scripts/first-deploy.sh
#
# 此脚本会：
# - 创建部署目录
# - 创建环境变量文件
# - 配置 systemd 服务
# - 配置 Nginx
# - 执行数据库迁移
# ============================================

set -e  # 遇到错误立即退出

# ============================================
# 加载配置
# ============================================
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 检查配置文件
if [ ! -f "${SCRIPT_DIR}/ops/config.sh" ]; then
    echo "❌ 配置文件不存在：${SCRIPT_DIR}/ops/config.sh"
    echo "请先创建配置文件，参考 scripts/ops/config.example.sh"
    exit 1
fi

source "${SCRIPT_DIR}/ops/config.sh"

# ============================================
# 颜色输出
# ============================================
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}🚀 开始首次部署 ${SYSTEM_NAME}...${NC}"
echo ""

# ============================================
# Step 1: 创建目录
# ============================================
echo -e "\n📁 [1/6] 创建部署目录..."

ssh ${SERVER} << ENDSSH
# 创建后端目录
sudo mkdir -p ${REMOTE_DIR}
sudo mkdir -p ${REMOTE_DIR}/logs
sudo mkdir -p ${REMOTE_DIR}-frontend

# 设置权限
sudo chown -R ubuntu:ubuntu ${REMOTE_DIR}
sudo chown -R ubuntu:ubuntu ${REMOTE_DIR}-frontend

echo "✅ 目录创建完成"
ls -la /var/www/ | grep ${SYSTEM_NAME}
ENDSSH

echo -e "${GREEN}✅ 目录创建完成${NC}"

# ============================================
# Step 2: 创建环境变量
# ============================================
echo -e "\n🔐 [2/6] 创建环境变量..."

echo "请在服务器上手动创建环境变量文件："
echo ""
echo "ssh ${SERVER} \"sudo nano ${REMOTE_DIR}/.env\""
echo ""
echo "粘贴以下内容（修改密码等敏感信息）："
echo ""
cat << 'EOF'
APP_ENV=production
APP_PORT=8084

# 数据库配置（通过 SSH 隧道）
DB_HOST=localhost
DB_PORT=5432
DB_USER=nexus_user
DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
DB_NAME=edit_business_db
DB_SSLMODE=disable

# Auth Center 配置
AUTH_CENTER_URL=https://os.crazyaigc.com
AUTH_CENTER_CALLBACK_URL=https://edit.crazyaigc.com/api/v1/auth/callback

# JWT 配置（生产环境必须修改！）
JWT_SECRET=change-this-secret-in-production-min-32-chars
JWT_ACCESS_TOKEN_EXPIRE=24h

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
EOF
echo ""
read -p "环境变量已创建？(y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "请先创建环境变量文件，然后重新运行此脚本"
    exit 1
fi

echo -e "${GREEN}✅ 环境变量配置完成${NC}"

# ============================================
# Step 3: 配置 systemd 服务
# ============================================
echo -e "\n⚙️ [3/6] 配置 systemd 服务..."

ssh ${SERVER} "sudo tee /etc/systemd/system/${SYSTEM_NAME}.service" < "${SCRIPT_DIR}/ops/systemd.template.service"

ssh ${SERVER} << ENDSSH
# 重载 systemd
sudo systemctl daemon-reload

# 启用服务（开机自启）
sudo systemctl enable ${SYSTEM_NAME}

echo "✅ systemd 服务配置完成"
sudo systemctl status ${SYSTEM_NAME} --no-pager --lines=5
ENDSSH

echo -e "${GREEN}✅ systemd 服务配置完成${NC}"

# ============================================
# Step 4: 配置 Nginx
# ============================================
echo -e "\n🌐 [4/6] 配置 Nginx..."

# 上传 Nginx 配置
scp "${SCRIPT_DIR}/ops/nginx.template.conf" ${SERVER}:/tmp/${SYSTEM_NAME}.conf

ssh ${SERVER} << ENDSSH
# 移动配置到 Nginx 目录
sudo mv /tmp/${SYSTEM_NAME}.conf /etc/nginx/sites-available/${SYSTEM_NAME}

# 创建软链接启用站点
sudo ln -sf /etc/nginx/sites-available/${SYSTEM_NAME} /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载 Nginx
sudo systemctl reload nginx

echo "✅ Nginx 配置完成"
ENDSSH

echo -e "${GREEN}✅ Nginx 配置完成${NC}"

# ============================================
# Step 5: 执行数据库迁移
# ============================================
echo -e "\n🗄️ [5/6] 执行数据库迁移..."

echo "请在服务器上执行数据库迁移："
echo ""
echo "ssh ${SERVER}"
echo "psql -h localhost -U nexus_user -d edit_business_db -f /var/www/${SYSTEM_NAME}/migrations/001_init_schema.up.sql"
echo ""
read -p "数据库迁移已执行？(y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "请先执行数据库迁移，然后重新运行此脚本"
    exit 1
fi

echo -e "${GREEN}✅ 数据库迁移完成${NC}"

# ============================================
# Step 6: 部署应用
# ============================================
echo -e "\n📦 [6/6] 部署应用..."

echo "现在可以执行部署脚本："
echo ""
echo "  ./scripts/deploy-production.sh"
echo ""

echo -e "${GREEN}✅ 首次部署配置完成！${NC}"
echo ""
echo "下一步："
echo "1. 执行部署脚本：./scripts/deploy-production.sh"
echo "2. 验证服务：ssh ${SERVER} \"sudo systemctl status ${SYSTEM_NAME}\""
echo "3. 访问应用：https://${DOMAIN}"
echo ""
