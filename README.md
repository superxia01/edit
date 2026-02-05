# Edit Business - 小红书数据采集系统

采集小红书笔记和博主信息的业务系统，通过 Chrome 插件实现数据自动化采集。

## 技术栈

### 前端
- **框架**: Vite 7+ + React 19+ + TypeScript 5+
- **表格**: TanStack Table v8
- **UI 组件**: shadcn/ui (Radix UI + Tailwind CSS 4.x)
- **状态管理**: Zustand
- **路由**: React Router 7+
- **HTTP 客户端**: Axios

### 后端
- **语言**: Go 1.21+
- **框架**: Gin (github.com/gin-gonic/gin)
- **ORM**: GORM (gorm.io/gorm)
- **数据库**: PostgreSQL 15
- **认证**: JWT + 账号中心集成

## 功能特性

### Chrome 插件
- 单篇笔记采集
- 博主笔记批量采集
- 博主信息采集
- 自动同步到业务系统

### Web 管理后台
- **数据概览**: 笔记数、博主数、互动数据统计
- **笔记管理**: 列表展示、筛选、查看详情、删除
- **博主管理**: 列表展示、粉丝数排序、删除
- **用户认证**: 通过账号中心统一登录

## 项目结构

```
edit-business/
├── frontend/                 # 前端项目
│   ├── src/
│   │   ├── components/       # UI 组件
│   │   ├── pages/            # 页面
│   │   ├── api/              # API 封装
│   │   ├── hooks/            # 自定义 Hooks
│   │   └── lib/              # 工具函数
│   ├── package.json
│   └── vite.config.ts
│
├── backend/                  # 后端项目
│   ├── cmd/server/           # 应用入口
│   ├── internal/
│   │   ├── handler/          # HTTP 处理器
│   │   ├── service/          # 业务逻辑
│   │   ├── repository/       # 数据访问
│   │   ├── model/            # 数据模型
│   │   ├── middleware/       # 中间件
│   │   ├── router/           # 路由配置
│   │   └── config/           # 配置管理
│   ├── pkg/database/         # 数据库连接
│   ├── migrations/           # 数据库迁移
│   ├── .env.example          # 环境变量模板
│   └── go.mod
│
├── scripts/                  # 部署脚本（不在 Git 中）
│   └── ops/                  # 运维配置（不在 Git 中）
│
├── DEPLOYMENT.md             # 部署文档
└── README.md                 # 本文档
```

## 快速开始

### 前端开发

```bash
cd frontend

# 安装依赖
npm install

# 本地开发
npm run dev

# 类型检查
npm run check

# 构建
npm run build
```

### 后端开发

```bash
cd backend

# 下载依赖
go mod download

# 本地运行
go run cmd/server/main.go

# 编译
go build -o edit-business ./cmd/server
```

## 环境配置

### 前端环境变量 (.env.production)

```bash
VITE_API_BASE_URL=https://edit.crazyaigc.com
```

### 后端环境变量 (.env)

参考 `.env.example` 文件，包含以下配置：

```bash
# 应用配置
APP_ENV=production
APP_PORT=8084

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=nexus_user
DB_PASSWORD=***
DB_NAME=edit_business_db

# JWT 配置
JWT_SECRET=***

# 认证中心
AUTH_CENTER_URL=https://os.crazyaigc.com
```

## 数据库迁移

```bash
# 执行迁移
psql -h localhost -U nexus_user -d edit_business_db -f migrations/001_init_schema.up.sql

# 回滚迁移
psql -h localhost -U nexus_user -d edit_business_db -f migrations/001_init_schema.down.sql
```

## API 接口

### 认证相关（无需认证）
- `POST /api/v1/auth/wechat` - 微信登录回调
- `GET /api/v1/auth/callback` - 认证回调

### 笔记相关
- `POST /api/v1/notes` - 同步单篇笔记（Chrome 插件）
- `POST /api/v1/notes/batch` - 批量同步笔记（Chrome 插件）
- `GET /api/v1/notes` - 分页查询笔记（需认证）
- `GET /api/v1/notes/:id` - 获取笔记详情（需认证）
- `DELETE /api/v1/notes/:id` - 删除笔记（需认证）

### 博主相关
- `POST /api/v1/bloggers` - 同步博主信息（Chrome 插件）
- `GET /api/v1/bloggers` - 分页查询博主（需认证）
- `DELETE /api/v1/bloggers/:id` - 删除博主（需认证）

### 统计数据
- `GET /api/v1/stats` - 获取统计数据（需认证）

## Chrome 插件

插件仓库：[xhs2feishu](https://github.com/keenchase/xhs2feishu)

### 安装
1. 下载插件源码
2. 修改 `api-config.js` 中的 `BASE_URL` 为本系统地址
3. 在 Chrome 中加载插件（开发者模式）

### 使用
1. 访问小红书笔记详情页，点击插件图标采集单篇笔记
2. 访问小红书博主主页，批量采集笔记列表
3. 采集博主基本信息

## 部署

详细部署步骤请参考 [DEPLOYMENT.md](./DEPLOYMENT.md)

### 本地构建

```bash
# 前端构建
cd frontend
npm run build

# 后端交叉编译
cd ../backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o edit-business ./cmd/server
```

### 部署脚本

使用部署脚本自动部署（需先配置 `scripts/ops` 目录）：

```bash
./scripts/deploy-production.sh
```

## 开发规范

本项目遵循 [KeenChase 技术规范 V3.0](https://github.com/keenchase/keenchase-standards)：

- 代码命名：PascalCase (Go/TS) / snake_case (DB)
- API 设计：RESTful + JSON 统一响应
- 数据库：UUID 主键 + PostgreSQL
- 认证：账号中心统一认证

## 故障排查

### 前端
- **构建失败**: 检查 Node.js 版本 (推荐 18+)
- **API 请求失败**: 检查 `VITE_API_BASE_URL` 配置

### 后端
- **编译失败**: 检查 Go 版本 (推荐 1.21+)
- **数据库连接失败**: 检查 SSH 隧道状态

### 数据库
- **连接超时**: 检查 SSH 隧道服务状态
- **权限错误**: 确认使用 `nexus_user` 用户

## License

MIT
