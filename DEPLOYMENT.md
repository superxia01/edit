# Edit Business 部署说明

本文档提供快速部署指南。详细的部署脚本使用说明请参考 [scripts/README.md](./scripts/README.md)。

## 快速部署

### 前置要求

- 本地已配置 SSH 别名（`~/.ssh/config`）
- 服务器已配置 SSH 隧道到数据库
- 本地已安装 Node.js 18+ 和 Go 1.21+

### 首次部署

```bash
# 1. 创建配置文件
cp scripts/ops/config.example.sh scripts/ops/config.sh
nano scripts/ops/config.sh

# 2. 执行首次部署
chmod +x scripts/first-deploy.sh
./scripts/first-deploy.sh
```

首次部署会自动完成：
- 创建部署目录
- 配置 systemd 服务
- 配置 Nginx
- 提示创建环境变量
- 提示执行数据库迁移

### 日常部署

```bash
# 执行部署脚本
chmod +x scripts/deploy-production.sh
./scripts/deploy-production.sh
```

## 目录结构

```
edit-business/
├── frontend/               # 前端项目
│   ├── dist/              # 构建输出
│   └── ...
├── backend/               # 后端项目
│   ├── edit-business      # 二进制文件
│   ├── migrations/        # 数据库迁移
│   └── ...
├── scripts/               # 部署脚本
│   ├── deploy-production.sh
│   ├── first-deploy.sh
│   └── ops/              # 运维配置（敏感）
└── DEPLOYMENT.md          # 本文档
```

## 环境变量

### 前端 (.env.production)

```bash
VITE_API_BASE_URL=https://edit.crazyaigc.com
```

### 后端 (.env)

```bash
APP_ENV=production
APP_PORT=8084

DB_HOST=localhost
DB_PORT=5432
DB_USER=nexus_user
DB_PASSWORD=***
DB_NAME=edit_business_db

AUTH_CENTER_URL=https://os.crazyaigc.com
JWT_SECRET=***
```

详细配置参考 `backend/.env.example`。

## 数据库迁移

```bash
# 执行迁移
psql -h localhost -U nexus_user -d edit_business_db -f backend/migrations/001_init_schema.up.sql

# 回滚迁移
psql -h localhost -U nexus_user -d edit_business_db -f backend/migrations/001_init_schema.down.sql
```

## 运维管理

### 查看服务状态

```bash
# 后端服务
ssh shanghai-tencent "sudo systemctl status edit-business"

# 查看日志
ssh shanghai-tencent "sudo journalctl -u edit-business -f"

# Nginx 状态
ssh shanghai-tencent "sudo systemctl status nginx"
```

### 回滚部署

```bash
# 登录服务器
ssh shanghai-tencent

# 查看备份
ls -la /var/www/edit-business/*.backup.*

# 回滚
sudo mv /var/www/edit-business/edit-business.backup.20260204_120000 /var/www/edit-business/edit-business
sudo systemctl restart edit-business
```

## 故障排查

### 服务启动失败

```bash
# 查看详细日志
ssh shanghai-tencent "sudo journalctl -u edit-business -n 50"

# 检查环境变量
ssh shanghai-tencent "sudo cat /var/www/edit-business/.env"

# 测试数据库连接
ssh shanghai-tencent "psql -h localhost -U nexus_user -d edit_business_db -c 'SELECT 1;'"
```

### 前端访问 404

```bash
# 检查前端文件
ssh shanghai-tencent "ls -la /var/www/edit-business-frontend/"

# 检查 Nginx 配置
ssh shanghai-tencent "sudo nginx -t"

# 查看 Nginx 日志
ssh shanghai-tencent "sudo tail -f /var/log/nginx/error.log"
```

## 更多信息

- **详细部署说明**: [scripts/README.md](./scripts/README.md)
- **开发文档**: [README.md](./README.md)
- **KeenChase 规范**: [keenchase-standards](https://github.com/keenchase/keenchase-standards)
