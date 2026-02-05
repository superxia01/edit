# 小红书数据采集系统 PRD

## 📚 技术规范

本项目遵循 **KeenChase 通用技术规范 V3.0**：

- **[系统架构与技术标准](https://github.com/keenchase/keenchase-standards/blob/main/architecture.md)**
- **[部署与服务管理](https://github.com/keenchase/keenchase-standards/blob/main/deployment-and-operations.md)**
- **[SSH 配置指南](https://github.com/keenchase/keenchase-standards/blob/main/ssh-setup.md)**
- **[数据库使用指南](https://github.com/keenchase/keenchase-standards/blob/main/database-guide.md)**
- **[安全规范](https://github.com/keenchase/keenchase-standards/blob/main/security.md)**
- **[API 接口说明](https://github.com/keenchase/keenchase-standards/blob/main/api.md)**

---

## 1. 项目概述

### 1.1 背景
将现有的 `xhs2feishu` Chrome 插件改造,使其不再同步到飞书/Coze/Keenchase等第三方平台,改为同步到自建的 edit-business 系统（自媒体创作工具的一部分）。

### 1.2 目标
- 保持插件采集功能不变
- 将数据同步目标改为 edit-business 系统
- 在 edit-business 系统中提供数据表格展示采集内容

### 1.3 技术栈（KeenChase V3.0 标准）

#### 前端
- **框架**: Vite 6+ + React 19+ + TypeScript 5+
- **表格**: TanStack Table (React Table v8)
- **UI 组件**: shadcn/ui (基于 Radix UI + Tailwind CSS)
- **状态管理**: Zustand / React Context
- **路由**: React Router 6+
- **HTTP 客户端**: Axios / Fetch API
- **构建工具**: Vite

#### 后端
- **语言**: Go 1.21+
- **框架**: Gin (github.com/gin-gonic/gin)
- **ORM**: GORM (gorm.io/gorm)
- **数据库**: PostgreSQL 15+ (杭州服务器统一数据库)
- **数据库驱动**: GORM PostgreSQL (gorm.io/driver/postgres)
- **认证**: JWT (github.com/golang-jwt/jwt/v5) + 账号中心集成
- **密码加密**: bcrypt (golang.org/x/crypto/bcrypt)
- **配置管理**: godotenv (github.com/joho/godotenv)
- **日志**: Zap (go.uber.org/zap)

#### 数据库
- **数据库**: PostgreSQL 15 (统一在杭州服务器)
- **用户**: nexus_user (统一数据库用户)
- **连接方式**: SSH 隧道 (localhost:5432)
- **隔离策略**: 独立数据库 edit_business_db

---

## 2. KeenChase V3.0 架构标准

### 2.1 系统架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              用户层 (浏览器)                                │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    edit-business 系统 (上海服务器)                          │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  前端 (Vite + React + TanStack Table + shadcn/ui)                     │   │
│  │  - Nginx 直接服务静态文件                                             │   │
│  │  - 域名: edit.crazyaigc.com                                           │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  后端 (Go + Gin)                                                      │   │
│  │  - Systemd 管理                                                       │   │
│  │  - 端口: 8084                                                         │   │
│  │  - SSH 隧道连接数据库                                                  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                 统一数据层 (杭州服务器 47.110.82.96)                         │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐    │
│  │         PostgreSQL 15 (端口5432, Docker)                            │    │
│  │                                                                     │    │
│  │  ┌──────────────────────────────────────────────────────────────┐ │    │
│  │  │  auth_center_db     (账号中心数据库)                           │ │    │
│  │  │  - users, user_accounts, sessions                              │ │    │
│  │  └──────────────────────────────────────────────────────────────┘ │    │
│  │                                                                     │    │
│  │  ┌──────────────────────────────────────────────────────────────┐ │    │
│  │  │  edit_business_db    (自媒体创作数据库)                        │ │    │
│  │  │  - notes, bloggers, users                                     │ │    │
│  │  └──────────────────────────────────────────────────────────────┘ │    │
│  └────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 服务器部署架构

#### 上海服务器 (101.35.120.199) - 应用服务器
- **用途**: 部署前后端应用
- **用户**: ubuntu (操作系统用户)
- **部署系统**:
  - edit.crazyaigc.com (自媒体创作系统，本项目)
  - os.crazyaigc.com (账号中心)
  - pr.crazyaigc.com (PR业务系统)
  - pixel.crazyaigc.com (AI生图系统)
  - quote.crazyaigc.com (报价系统)

#### 杭州服务器 (47.110.82.96) - 统一数据库服务器
- **用途**: 统一数据存储中心
- **用户**: root (操作系统用户)
- **数据库用户**: nexus_user (PostgreSQL 超级用户)
- **数据库密码**: hRJ9NSJApfeyFDraaDgkYowY
- **数据库列表**:
  - auth_center_db (账号中心)
  - edit_business_db (自媒体创作，本项目)
  - pr_business_db (PR业务)
  - pixel_business_db (AI生图)
  - quote_business_db (报价系统)

---

## 3. 功能需求

### 3.1 Chrome 插件功能(保持不变)

插件支持三种采集模式:

#### 3.1.1 单篇笔记采集
在笔记详情页采集单条笔记数据:

| 字段 | 类型 | 说明 |
|------|------|------|
| url | string | 笔记链接 |
| title | string | 标题 |
| author | string | 作者昵称 |
| content | string | 正文内容 |
| tags | string[] | 话题标签(去除#号) |
| imageUrls | string[] | 所有图片URL |
| videoUrl | string | 视频URL(可选) |
| noteType | string | 笔记类型(图文/视频) |
| coverImageUrl | string | 封面图片URL |
| likes | number | 点赞数 |
| collects | number | 收藏数 |
| comments | number | 评论数 |
| publishDate | timestamp | 发布时间 |
| captureTimestamp | timestamp | 采集时间戳 |

#### 3.1.2 博主笔记批量采集
在博主主页批量采集笔记列表:

| 字段 | 类型 | 说明 |
|------|------|------|
| title | string | 标题 |
| url | string | 笔记链接 |
| author | string | 作者昵称 |
| likes | number | 点赞数 |
| image | string | 封面图片URL |

#### 3.1.3 博主信息采集
采集博主基本信息:

| 字段 | 类型 | 说明 |
|------|------|------|
| bloggerName | string | 博主名称 |
| avatarUrl | string | 头像链接 |
| bloggerId | string | 小红书号 |
| description | string | 简介 |
| followersCount | number | 粉丝数 |
| bloggerUrl | string | 博主主页链接 |
| captureTimestamp | timestamp | 采集时间戳 |

---

## 4. 后端设计（遵循 KeenChase V3.0 标准）

### 4.1 Go 项目结构

```
edit-business/                      # ✅ 后端已基本完成
├── cmd/
│   └── server/
│       └── main.go              # ✅ 程序入口
├── internal/
│   ├── handler/
│   │   ├── note.go              # ✅ 笔记 HTTP 处理器
│   │   ├── blogger.go           # ✅ 博主 HTTP 处理器
│   │   ├── auth.go              # ✅ 认证相关（账号中心集成）
│   │   ├── user.go              # ✅ 用户 HTTP 处理器
│   │   └── response.go          # ✅ 统一响应格式
│   ├── middleware/
│   │   └── jwt.go               # ✅ JWT 认证中间件
│   ├── model/
│   │   ├── note.go              # ✅ 笔记 GORM 模型
│   │   ├── blogger.go           # ✅ 博主 GORM 模型
│   │   ├── user.go              # ✅ 用户 GORM 模型（账号中心关联）
│   │   └── types.go             # ✅ 自定义类型
│   ├── repository/
│   │   ├── note.go              # ✅ 笔记数据库操作
│   │   ├── blogger.go           # ✅ 博主数据库操作
│   │   └── user.go              # ✅ 用户数据库操作
│   ├── service/
│   │   ├── note.go              # ✅ 笔记业务逻辑
│   │   ├── blogger.go           # ✅ 博主业务逻辑
│   │   └── user.go              # ✅ 用户业务逻辑
│   ├── router/
│   │   └── router.go            # ✅ 路由配置（含 CORS、JWT）
│   └── config/
│       └── config.go            # ✅ 配置管理
├── pkg/
│   └── database/
│       └── postgres.go          # ✅ 数据库连接（SSH 隧道）
├── .env                         # ✅ 环境变量
├── go.mod                       # ✅
├── go.sum                       # ✅
├── Dockerfile                   # ⏳ 待实现
└── deployments/
    └── production.sh            # ⏳ 待实现
```

### 4.2 数据库设计（遵循 KeenChase 命名规范）

#### 4.2.1 核心规范

**⚠️ 强制规则**：
- ✅ **表名**: `snake_case`，复数形式
- ✅ **列名**: `snake_case`，全部小写
- ✅ **主键**: UUID (不是 Auto Increment INT)
- ✅ **外键**: `{table}_{column}_fkey`
- ✅ **索引**: `{table}_{column}_idx`
- ✅ **时间戳**: `{column}_at` (timestamp with time zone)
- ✅ **JSON 字段**: JSONB (不是 JSON)

#### 4.2.2 notes 表（笔记数据）

```sql
-- ✅ 正确：使用 UUID 主键 + snake_case 命名
CREATE TABLE notes (
  -- 主键：强制使用 UUID
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- 业务字段：snake_case 命名
  url VARCHAR(500) UNIQUE NOT NULL,
  title VARCHAR(500),
  author VARCHAR(100),
  content TEXT,
  tags TEXT[],                    -- PostgreSQL array type
  image_urls TEXT[],              -- PostgreSQL array type
  video_url VARCHAR(500),
  note_type VARCHAR(20),          -- '图文' or '视频'
  cover_image_url VARCHAR(500),
  likes INTEGER DEFAULT 0,
  collects INTEGER DEFAULT 0,
  comments INTEGER DEFAULT 0,

  -- 时间戳：timestamp with time zone
  publish_date BIGINT,
  capture_timestamp BIGINT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ✅ 索引命名：{table}_{column}_idx
CREATE INDEX idx_notes_author ON notes(author);
CREATE INDEX idx_notes_publish_date ON notes(publish_date DESC);
CREATE INDEX idx_notes_capture_timestamp ON notes(capture_timestamp DESC);
CREATE INDEX idx_notes_tags ON notes USING GIN(tags);
CREATE INDEX idx_notes_note_type ON notes(note_type);
```

#### 4.2.3 bloggers 表（博主信息）

```sql
CREATE TABLE bloggers (
  -- 主键：UUID
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- 业务字段
  xhs_id VARCHAR(50) UNIQUE,        -- 小红书号
  blogger_name VARCHAR(100),
  avatar_url VARCHAR(500),
  description TEXT,
  followers_count INTEGER DEFAULT 0,
  blogger_url VARCHAR(500),
  capture_timestamp BIGINT NOT NULL,

  -- 时间戳
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_bloggers_xhs_id ON bloggers(xhs_id);
CREATE INDEX idx_bloggers_followers ON bloggers(followers_count DESC);
```

#### 4.2.4 users 表（关联账号中心）

```sql
CREATE TABLE users (
  -- 本地主键
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- ✅ 关联账号中心（强制）
  auth_center_user_id UUID UNIQUE NOT NULL,

  -- 业务字段
  role VARCHAR(50) DEFAULT 'USER',
  profile JSONB,

  -- 时间戳
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

  -- 外键约束
  CONSTRAINT users_auth_center_user_id_fkey
    FOREIGN KEY (auth_center_user_id)
    REFERENCES auth_center_db.users(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);

-- 索引
CREATE INDEX idx_users_auth_center_user_id ON users(auth_center_user_id);
```

### 4.3 GORM 模型示例（遵循 KeenChase 命名规范）

```go
// ✅ 正确：结构体 PascalCase，JSON camelCase，GORM column snake_case
type Note struct {
    ID              UUID              `gorm:"primaryKey;column:id;type:uuid" json:"id"`
    URL             string            `gorm:"uniqueIndex;column:url;type:varchar(500)" json:"url"`
    Title           string            `gorm:"column:title;type:varchar(500)" json:"title"`
    Author          string            `gorm:"column:author;type:varchar(100)" json:"author"`
    Content         string            `gorm:"column:content;type:text" json:"content"`
    Tags            pq.StringArray    `gorm:"column:tags;type:text[]" json:"tags"`
    ImageURLs       pq.StringArray    `gorm:"column:image_urls;type:text[]" json:"imageUrls"`
    VideoURL        *string           `gorm:"column:video_url;type:varchar(500)" json:"videoUrl,omitempty"`
    NoteType        string            `gorm:"column:note_type;type:varchar(20)" json:"noteType"`
    CoverImageURL   string            `gorm:"column:cover_image_url;type:varchar(500)" json:"coverImageUrl"`
    Likes           int32             `gorm:"column:likes;type:integer" json:"likes"`
    Collects        int32             `gorm:"column:collects;type:integer" json:"collects"`
    Comments        int32             `gorm:"column:comments;type:integer" json:"comments"`
    PublishDate     int64             `gorm:"column:publish_date;type:bigint" json:"publishDate"`
    CaptureTimestamp int64            `gorm:"column:capture_timestamp;type:bigint" json:"captureTimestamp"`
    CreatedAt       time.Time         `gorm:"column:created_at;type:timestamp with time zone" json:"createdAt"`
    UpdatedAt       time.Time         `gorm:"column:updated_at;type:timestamp with time zone" json:"updatedAt"`
}

// ✅ 指定表名（复数 + snake_case）
func (Note) TableName() string {
    return "notes"
}
```

### 4.4 API 接口设计（RESTful 标准）

#### 4.4.1 基础规范

- ✅ 使用名词复数: `/api/v1/notes`, `/api/v1/bloggers`
- ✅ HTTP 方法语义化:
  - `GET` - 查询
  - `POST` - 创建
  - `PUT` - 完整更新
  - `DELETE` - 删除

#### 4.4.2 接口列表

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/v1/notes/single | 同步单篇笔记 | ❌ |
| POST | /api/v1/notes/batch | 批量同步笔记 | ❌ |
| POST | /api/v1/bloggers | 同步博主信息 | ❌ |
| GET | /api/v1/notes | 分页查询笔记 | ✅ |
| GET | /api/v1/notes/:id | 获取笔记详情 | ✅ |
| DELETE | /api/v1/notes/:id | 删除笔记 | ✅ |
| GET | /api/v1/bloggers | 分页查询博主 | ✅ |
| DELETE | /api/v1/bloggers/:id | 删除博主 | ✅ |
| GET | /api/v1/stats | 统计数据 | ✅ |

**说明**：
- 插件同步接口（POST）无需认证，允许 Chrome 插件直接调用
- 查询/删除接口需要 JWT Token 认证

#### 4.4.3 统一响应格式

**成功响应**：
```json
{
  "success": true,
  "data": {
    "id": "uuid-xxx",
    "title": "笔记标题"
  }
}
```

**列表响应**：
```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "total": 100
  }
}
```

**错误响应**：
```json
{
  "success": false,
  "error": "错误信息（用户可读）",
  "errorCode": "NOTE_NOT_FOUND"
}
```

### 4.5 数据库连接配置（SSH 隧道）

#### 4.5.1 连接方式

**⚠️ 重要：通过 SSH 隧道连接数据库**

```
上海服务器 (Go 应用)
  └─ SSH 隧道 (localhost:5432 → 47.110.82.96:5432)
      └─ 杭州服务器 (PostgreSQL)
```

#### 4.5.2 环境变量配置

**代码仓库** (`.env.example`):
```bash
# ============================================
# 应用配置
# ============================================
APP_ENV=production
APP_PORT=8084
APP_NAME=edit-business
APP_DEBUG=false

# ============================================
# 数据库配置（通过 SSH 隧道）
# ============================================
# ⚠️ 重要：所有系统统一使用 nexus_user
# ⚠️ 必须使用 localhost（通过 SSH 隧道转发）
DB_HOST=localhost
DB_PORT=5432
DB_USER=nexus_user
DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
DB_NAME=edit_business_db
DB_SSLMODE=disable

# ============================================
# Auth Center 配置
# ============================================
AUTH_CENTER_URL=https://os.crazyaigc.com
AUTH_CENTER_CALLBACK_URL=https://edit.crazyaigc.com/api/v1/auth/callback

# ============================================
# 前端地址
# ============================================
FRONTEND_URL=https://edit.crazyaigc.com

# ============================================
# CORS 白名单（插件域名）
# ============================================
ALLOWED_ORIGINS=https://edit.crazyaigc.com,chrome-extension://*

# ============================================
# JWT 配置
# ============================================
JWT_SECRET={CHANGE_THIS_IN_PRODUCTION}
JWT_ACCESS_TOKEN_EXPIRE=24h

# ============================================
# 日志配置
# ============================================
LOG_LEVEL=info
LOG_FORMAT=json
```

#### 4.5.3 数据库初始化代码

```go
// pkg/database/postgres.go
package database

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewPostgresDB() (*gorm.DB, error) {
    // ⚠️ 通过 SSH 隧道连接：localhost:5432
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        os.Getenv("DB_HOST"),      // localhost
        os.Getenv("DB_PORT"),      // 5432
        os.Getenv("DB_USER"),      // nexus_user
        os.Getenv("DB_PASSWORD"),  // hRJ9NSJApfeyFDraaDgkYowY
        os.Getenv("DB_NAME"),      // edit_business_db
        os.Getenv("DB_SSLMODE"),   // disable
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })

    if err != nil {
        return nil, err
    }

    log.Println("✅ 数据库连接成功（通过 SSH 隧道）")
    return db, nil
}
```

---

## 5. 前端设计（KeenChase V3.0 标准）

### 5.1 技术架构

#### 核心技术栈

```
前端架构 (Vite + React + TanStack Table + shadcn/ui)
├── Vite 6+                # 构建工具
├── React 19+               # UI 框架
├── TypeScript 5+           # 类型系统
├── TanStack Table v8       # 表格逻辑层 (Headless)
├── shadcn/ui              # UI 组件库 (渲染层)
│   ├── Table              # 基础表格组件
│   ├── Button             # 按钮
│   ├── Input              # 输入框
│   ├── Badge              # 标签
│   ├── Dialog             # 对话框
│   └── Dropdown           # 下拉菜单
├── Tailwind CSS           # 样式框架
└── Zustand                # 状态管理（可选）
```

### 5.2 前端项目结构

```
edit-business/
├── frontend/                      # 前端目录 ✅ 已完成
│   ├── src/
│   │   ├── components/
│   │   │   ├── Navigation.tsx       # ✅ 导航组件
│   │   │   └── ui/                # shadcn/ui 组件
│   │   │       ├── table.tsx      # ✅
│   │   │       ├── button.tsx     # ✅
│   │   │       ├── input.tsx      # ✅
│   │   │       ├── badge.tsx      # ✅
│   │   │       ├── card.tsx       # ✅
│   │   │       └── data-table.tsx # ✅ TanStack Table 封装
│   │   ├── pages/
│   │   │   ├── NotesListPage.tsx      # ✅ 笔记列表页
│   │   │   ├── BloggersListPage.tsx   # ✅ 博主列表页
│   │   │   ├── Dashboard.tsx        # ⏳ 仪表板（待实现）
│   │   │   └── Login.tsx            # ⏳ 登录页（待实现）
│   │   ├── lib/
│   │   │   ├── api.ts            # ✅ API 请求封装
│   │   │   └── utils.ts          # ✅ 工具函数
│   │   ├── App.tsx               # ✅
│   │   └── main.tsx              # ✅
│   ├── package.json              # ✅
│   ├── vite.config.ts            # ✅
│   └── tailwind.config.js        # ✅
└── backend/                      # 后端目录 ✅
```

### 5.3 DataTable 组件设计

**组件层次**:
```
DataTable (封装组件)
  ├── TanStack Table (逻辑层)
  │   ├── 排序
  │   ├── 过滤
  │   ├── 分页
  │   └── 虚拟滚动(可选)
  └── shadcn/ui Table (渲染层)
      ├── Table
      ├── TableHeader
      ├── TableBody
      ├── TableCell
      └── TableRow
```

### 5.4 笔记列表页列定义

```typescript
// src/pages/NotesTable.tsx
import { ColumnDef } from '@tanstack/react-table'
import { Badge } from '@/components/ui/badge'

export const columns: ColumnDef<Note>[] = [
  // 封面图
  {
    accessorKey: 'coverImageUrl',
    header: '封面',
    cell: ({ row }) => (
      <img src={row.getValue('coverImageUrl')} className="w-12 h-12 rounded" />
    ),
  },

  // 标题(可点击跳转)
  {
    accessorKey: 'title',
    header: '标题',
    cell: ({ row }) => (
      <a href={row.original.url} target="_blank" className="text-blue-600">
        {row.getValue('title')}
      </a>
    ),
  },

  // 作者
  {
    accessorKey: 'author',
    header: '作者',
  },

  // 类型(Badge)
  {
    accessorKey: 'noteType',
    header: '类型',
    cell: ({ row }) => (
      <Badge variant={row.getValue('noteType') === '视频' ? 'default' : 'secondary'}>
        {row.getValue('noteType')}
      </Badge>
    ),
  },

  // 标签(多个 Badge)
  {
    accessorKey: 'tags',
    header: '标签',
    cell: ({ row }) => {
      const tags = row.getValue('tags') as string[]
      return (
        <div className="flex gap-1">
          {tags.slice(0, 2).map(tag => (
            <Badge key={tag} variant="outline">#{tag}</Badge>
          ))}
          {tags.length > 2 && <Badge>+{tags.length - 2}</Badge>}
        </div>
      )
    },
  },

  // 互动数据
  {
    id: 'interaction',
    header: '互动',
    cell: ({ row }) => (
      <div className="flex gap-2 text-sm">
        <span>👍 {row.original.likes}</span>
        <span>⭐ {row.original.collects}</span>
        <span>💬 {row.original.comments}</span>
      </div>
    ),
  },

  // 发布时间(可排序)
  {
    accessorKey: 'publishDate',
    header: '发布时间',
    cell: ({ row }) => formatDate(row.getValue('publishDate')),
  },

  // 操作按钮
  {
    id: 'actions',
    header: '操作',
    cell: ({ row }) => (
      <DropdownMenu>
        <DropdownMenuItem onClick={() => handleView(row.original)}>查看</DropdownMenuItem>
        <DropdownMenuItem onClick={() => handleDelete(row.original.id)}>删除</DropdownMenuItem>
      </DropdownMenu>
    ),
  },
]
```

---

## 6. Chrome 插件改造 ✅ 已完成

### 6.1 修改文件清单

| 文件 | 改造内容 | 状态 |
|------|----------|------|
| `api-config.js` | 修改 BASE_URL 为 edit-business 后端地址 | ✅ |
| `sidebar.js` | 修改同步函数,调用新的 API 接口,移除用户验证逻辑 | ✅ |
| `sidebar.html` | 移除设置标签页，简化 UI | ✅ |
| `manifest.json` | 添加 edit-business 域名到 host_permissions | ✅ |

### 6.2 具体改造内容

#### Step 1: 修改 api-config.js ✅

```javascript
const API_CONFIG = {
    // 本地开发
    BASE_URL: 'http://localhost:8084',

    // 生产环境
    // BASE_URL: 'https://edit.crazyaigc.com',

    ENDPOINTS: {
        SYNC_SINGLE_NOTE: '/api/v1/notes',
        SYNC_BLOGGER_NOTES: '/api/v1/notes/batch',
        SYNC_BLOGGER_INFO: '/api/v1/bloggers',
    }
};
```

#### Step 2: 简化 sidebar.js ✅

**移除内容**:
- ✅ `verifyUserOrder()` 函数及相关用户验证逻辑
- ✅ 飞书配置输入框 (ordeid, basetoken, knowledgeurl, bloggerurl, blogger_noteurl)
- ✅ 订单验证相关 UI
- ✅ 配置标签页（configTab）
- ✅ 下载媒体功能

**修改同步函数**:
- ✅ `syncSingleNote()` → POST `/api/v1/notes`
- ✅ `syncBatchNotes()` → POST `/api/v1/notes/batch`
- ✅ `syncBloggerInfo()` → POST `/api/v1/bloggers`

**代码精简**: 从 1,739 行 → 1,052 行（减少 39.5%）

#### Step 3: 修改 manifest.json ✅

```json
{
  "manifest_version": 3,
  "name": "edit-business-crawler",
  "version": "2.0.0",
  "description": "小红书数据采集并同步到 edit-business 系统",
  "permissions": [
    "activeTab",
    "scripting",
    "storage"
  ],
  "host_permissions": [
    "*://www.xiaohongshu.com/*",
    "*://localhost:8084/*",
    "*://edit.crazyaigc.com/*"
  ],
  "action": {
    "default_icon": {
      "16": "images/icon16.png",
      "48": "images/icon48.png",
      "128": "images/icon128.png"
    }
  },
  "icons": {
    "16": "images/icon16.png",
    "48": "images/icon48.png",
    "128": "images/icon128.png"
  },
  "content_scripts": [{
    "matches": ["*://www.xiaohongshu.com/*"],
    "js": ["content.js"]
  }],
  "background": {
    "service_worker": "background.js"
  },
  "side_panel": {
    "default_path": "sidebar.html"
  }
}
```

---

## 7. 部署规范（KeenChase V3.0）

### 7.1 目录命名规范

```
/var/www/
├── edit-business/                # 后端目录
│   ├── edit-business-api          # 可执行文件
│   ├── .env                       # 环境变量（服务器独立）
│   ├── .env.example               # 环境变量模板
│   └── logs/                      # 日志目录
│
└── edit-business-frontend/        # 前端目录
    ├── index.html
    └── assets/
```

### 7.2 环境变量管理

**⚠️ 核心原则：环境变量与代码分离，部署不覆盖配置**

**代码仓库**:
```
backend/
├── .env.example          # ✅ 环境变量模板（提交到 Git）
└── .env.local            # 本地开发（不提交）
```

**服务器上**:
```
/var/www/edit-business/
├── edit-business-api
├── .env                   # ✅ 实际环境变量（首次手动创建）
└── .env.backup            # 自动备份
```

**⚠️ 部署时不覆盖 .env 文件**

### 7.3 部署脚本

```bash
#!/bin/bash
# deployments/production.sh

set -e

SYSTEM_NAME="edit-business"
BINARY_NAME="edit-business-api"
DOMAIN="edit.crazyaigc.com"
SERVER="shanghai-tencent"
REMOTE_DIR="/var/www/${SYSTEM_NAME}"

echo "🚀 开始部署 ${SYSTEM_NAME}..."

# ============================================
# 前端部署
# ============================================
echo "📦 [1/3] 部署前端..."

cd frontend
npm run build

rsync -avz --delete \
  --exclude '*.map' \
  dist/ \
  ${SERVER}:${REMOTE_DIR}-frontend/

echo "✅ 前端部署完成"

# ============================================
# 后端部署
# ============================================
echo "📦 [2/3] 部署后端..."

cd ../backend

# 交叉编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -o ${BINARY_NAME} \
  cmd/server/main.go

# 上传二进制
scp ${BINARY_NAME} ${SERVER}:${REMOTE_DIR}/

# 重启服务（不覆盖 .env）
ssh ${SERVER} << ENDSSH
cd ${REMOTE_DIR}

# 备份旧二进制
if [ -f ${BINARY_NAME} ]; then
  mv ${BINARY_NAME} ${BINARY_NAME}.backup.\$(date +%Y%m%d_%H%M%S)
fi

# 重命名新二进制
mv ${BINARY_NAME}-new ${BINARY_NAME}

# 重启服务
sudo systemctl restart ${SYSTEM_NAME}

# 等待启动
sleep 3

# 检查状态
sudo systemctl status ${SYSTEM_NAME} --no-pager
ENDSSH

echo "✅ 后端部署完成"

# ============================================
# 验证部署
# ============================================
echo "🔍 [3/3] 验证部署..."

sleep 2
curl -f https://${DOMAIN}/health || echo "⚠️ 健康检查失败"

echo ""
echo "🎉 部署完成！"
echo ""
echo "📍 访问地址："
echo "  前端: https://${DOMAIN}"
```

---

## 8. 账号中心集成

### 8.1 认证流程

```
用户在 edit.crazyaigc.com
  ↓ 点击"微信登录"
  → 前端跳转到 /api/auth/wechat/login
  ↓
后端重定向到账号中心
  → 用户短暂看到 os.crazyaigc.com
  ↓
账号中心重定向到微信授权页面
  → 用户扫码/授权
  ↓
账号中心回调到 edit.crazyaigc.com/api/auth/callback
  ↓
业务系统后端接收 userId + token
  → 验证 token
  → 创建/获取本地用户
  → 设置 session
  → 跳转到 /dashboard
  ↓
用户登录完成 ✅
```

### 8.2 用户表设计

```sql
CREATE TABLE users (
  -- 本地主键
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- ✅ 关联账号中心（强制）
  auth_center_user_id UUID UNIQUE NOT NULL,

  -- 业务字段
  role VARCHAR(50) DEFAULT 'USER',
  profile JSONB,

  -- 时间戳
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

  -- 外键约束
  CONSTRAINT users_auth_center_user_id_fkey
    FOREIGN KEY (auth_center_user_id)
    REFERENCES auth_center_db.users(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);

-- 索引
CREATE INDEX idx_users_auth_center_user_id ON users(auth_center_user_id);
```

### 8.3 CORS 配置（插件白名单）

```go
// 允许 Chrome 插件调用
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
    allowedOrigins := []string{
        "https://edit.crazyaigc.com",
        "chrome-extension://*",  // ✅ 允许所有 Chrome 插件
    }

    originMap := make(map[string]bool)
    for _, origin := range allowedOrigins {
        originMap[origin] = true
    }

    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        if origin != "" {
            if !originMap[origin] && !strings.HasPrefix(origin, "chrome-extension://") {
                c.JSON(403, gin.H{
                    "success": false,
                    "error":   "域名未在白名单中",
                })
                c.Abort()
                return
            }
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
        }

        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

---

## 9. 开发排期

| 阶段 | 任务 | 预计时间 |
|------|------|----------|
| **Phase 1** | 前端基础搭建 | 1 天 |
| | - 初始化 Vite + React 项目 | |
| | - 安装 TanStack Table + shadcn/ui | |
| | - 创建 DataTable 组件 | |
| **Phase 2** | 后端开发 | 3-4 天 |
| | - 搭建 Go 项目框架 | |
| | - 数据库设计与迁移 | |
| | - 实现 CRUD API | |
| | - CORS/认证中间件 | |
| **Phase 3** | 前端页面开发 | 2-3 天 |
| | - 笔记列表页 | |
| | - 博主列表页 | |
| | - API 请求封装 | |
| | - 登录集成（账号中心） | |
| **Phase 4** | 插件改造 | 1 天 |
| | - 修改 API 配置 | |
| | - 简化同步逻辑 | |
| **Phase 5** | 联调测试 | 2 天 |
| | - 前后端联调 | |
| | - 插件与后端联调 | |
| **Phase 6** | 部署上线 | 1 天 |

**总计**: 约 10-12 天

---

## 10. 安全规范（KeenChase V3.0）

### 10.1 JWT Token 规范

```go
// ✅ Token 结构
type Claims struct {
    UserID string `json:"userId"`
    jwt.RegisteredClaims
}

// ✅ 标准配置
- 算法: HS256
- 有效期: 24小时
- 签名密钥: 最少32字符
- 存储: Cookie (httpOnly, secure, sameSite)
```

### 10.2 密码安全

```go
// ✅ 密码哈希（强制 bcrypt）
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 10.3 数据库安全

**连接字符串规范**:
```bash
# ✅ 通过 SSH 隧道（已加密，数据库层可 disable）
DATABASE_URL="postgresql://nexus_user:hRJ9NSJApfeyFDraaDgkYowY@localhost:5432/edit_business_db?sslmode=disable"

# ❌ 禁止直连且不使用SSL
DATABASE_URL="postgresql://nexus_user:pass@47.110.82.96:5432/db?sslmode=disable"
```

**SQL 注入防护**:
```go
// ✅ 使用 GORM 参数化查询
db.Where("id = ?", noteID).First(&note)

// ❌ 禁止字符串拼接
db.Where("id = '" + noteID + "'").First(&note)
```

---

## 11. 附录

### 11.1 数据流

```
用户操作(点击采集)
    ↓
插件 content.js (提取数据)
    ↓
插件 sidebar.js (组装数据)
    ↓
HTTP POST → Go API (上海服务器)
    ↓
Go Handler → Service → Repository
    ↓
PostgreSQL (杭州服务器，通过 SSH 隧道)
    ↓
HTTP GET ← React 前端
    ↓
TanStack Table (处理数据)
    ↓
shadcn/ui Table (渲染)
    ↓
用户查看
```

### 11.2 关键技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 表格方案 | TanStack Table | Headless、性能好、灵活 |
| UI 组件 | shadcn/ui | 设计一致、代码可控 |
| 样式方案 | Tailwind CSS | 快速开发、响应式 |
| 后端语言 | Go | 高性能、易部署 |
| 数据库 | PostgreSQL | Array 类型、GIN 索引 |
| 认证方式 | 账号中心统一认证 | 用户体系统一 |

### 11.3 参考资源

**KeenChase 技术规范**:
- [架构标准](https://github.com/keenchase/keenchase-standards/blob/main/architecture.md)
- [部署规范](https://github.com/keenchase/keenchase-standards/blob/main/deployment-and-operations.md)
- [数据库指南](https://github.com/keenchase/keenchase-standards/blob/main/database-guide.md)
- [安全规范](https://github.com/keenchase/keenchase-standards/blob/main/security.md)
- [API 文档](https://github.com/keenchase/keenchase-standards/blob/main/api.md)

**技术文档**:
- [TanStack Table 文档](https://tanstack.com/table/latest)
- [shadcn/ui 文档](https://ui.shadcn.com/)
- [Gin 框架](https://gin-gonic.com/)
- [GORM 文档](https://gorm.io/docs/)
- [PostgreSQL 文档](https://www.postgresql.org/docs/)

### 11.4 常用命令

```bash
# === SSH 连接 ===
ssh shanghai-tencent      # 上海服务器
ssh hangzhou-ali          # 杭州数据库服务器

# === 服务管理 ===
sudo systemctl status edit-business
sudo systemctl restart edit-business
sudo journalctl -u edit-business -f

# === Nginx 管理 ===
sudo nginx -t
sudo systemctl reload nginx

# === 数据库连接（通过 SSH 隧道）===
psql -h localhost -p 5432 -U nexus_user -d edit_business_db
```

---

## 12. 风险与挑战

### 12.1 技术风险

| 风险 | 影响 | 应对措施 |
|------|------|----------|
| 数据量过大导致前端卡顿 | 高 | 启用虚拟滚动,服务端分页 |
| 小红书页面结构变化 | 中 | 定期维护选择器,增强容错 |
| SSH 隧道中断 | 中 | 自动重连机制，监控脚本 |
| CORS 跨域问题 | 低 | 后端配置 CORS 白名单 |

### 12.2 业务风险

| 风险 | 影响 | 应对措施 |
|------|------|----------|
| 采集频率限制 | 中 | 添加限流,随机延迟 |
| 数据存储成本 | 中 | 定期清理旧数据,图片压缩 |
| 账号中心依赖 | 低 | 本地用户表做缓存 |

---

**文档版本**: v3.1 (KeenChase V3.0 标准 + 开发进度追踪)
**创建日期**: 2026-02-04
**最后更新**: 2026-02-04
**v3.0 更新内容**:
- ✅ 完全遵循 KeenChase 技术规范 V3.0
- ✅ 数据库使用 UUID 主键
- ✅ 通过 SSH 隧道连接数据库
- ✅ 统一使用 nexus_user 数据库用户
- ✅ 集成账号中心认证
- ✅ 标准化部署规范
- ✅ CORS 白名单配置（支持 Chrome 插件）

**v3.1 开发进度**:
### ✅ 已完成 (核心功能)

#### 后端 (Go + Gin)
- ✅ 项目结构搭建
- ✅ GORM 数据模型
- ✅ 数据库连接（SSH 隧道）
- ✅ Repository 层（CRUD 操作）
- ✅ Service 层（业务逻辑）
- ✅ Handler 层（HTTP 接口）
- ✅ JWT 认证中间件
- ✅ 账号中心集成
- ✅ 路由配置（含 CORS、JWT）

#### 前端 (Vite + React + TypeScript)
- ✅ 项目初始化（Vite + React + TypeScript）
- ✅ TanStack Table + shadcn/ui 集成
- ✅ DataTable 组件（排序、分页）
- ✅ API 客户端
- ✅ 笔记列表页面
- ✅ 博主列表页面
- ✅ 导航栏组件
- ✅ React Router 配置

#### Chrome 插件
- ✅ manifest.json 修改
- ✅ api-config.js 配置更新
- ✅ sidebar.html 简化（移除设置标签页）
- ✅ sidebar.js 重构（移除飞书集成、订单验证）
  - 从 1,739 行精简到 1,052 行（减少 39.5%）

### ⏳ 待实现

#### 前端
- ⏳ 登录页面（账号中心 OAuth）
- ⏳ 仪表板页面
- ⏳ 自定义 Hooks（useAuth, useNotes, useBloggers）
- ⏳ 筛选功能（按作者、标签筛选）
- ⏳ 笔记详情查看
- ⏳ 删除确认对话框
- ⏳ 加载状态 & 错误处理优化

#### 后端
- ⏳ 响应格式统一（当前与 PRD 不完全一致）
- ⏳ 日志系统（Zap 集成）
- ⏳ .env.example 环境变量模板
- ⏳ Dockerfile
- ⏳ Systemd 服务配置

#### 运维配置
- ⏳ SSH 隧道配置
- ⏳ Systemd 服务文件
- ⏳ Nginx 反向代理配置
- ⏳ 数据库迁移脚本执行

#### 测试
- ⏳ 前后端联调
- ⏳ Chrome 插件与后端联调
- ⏳ 采集功能测试

