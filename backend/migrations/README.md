# 数据库迁移指南

## 概述

本次迁移将用户信息从 JSONB 字段迁移到独立字段，提升查询性能并避免 JSONB 解析开销。

## 迁移时间

- **日期**: 2026-02-06
- **版本**: v1.0.0

## 迁移内容

### 新增字段

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| `union_id` | VARCHAR(255) | UNIQUE | 微信 UnionID |
| `nickname` | VARCHAR(100) | - | 用户昵称 |
| `avatar_url` | VARCHAR(500) | - | 用户头像URL |
| `phone_number` | VARCHAR(255) | UNIQUE | 手机号 |
| `email` | VARCHAR(255) | UNIQUE | 邮箱 |

### 数据来源

所有新字段的数据都从现有的 `profile` JSONB 字段迁移：

- `union_id` ← `profile->>'unionId'`
- `nickname` ← `profile->>'nickname'`
- `avatar_url` ← `profile->>'avatarUrl'`
- `phone_number` ← `profile->>'phoneNumber'`
- `email` ← `profile->>'email'`

## 执行步骤

### 1. 备份数据库（重要！）

```bash
# 在服务器上执行
ssh shanghai-tencent
pg_dump -U nexus_user -h localhost edit_business_db > backup_before_migration_$(date +%Y%m%d_%H%M%S).sql
```

### 2. 执行迁移脚本

```bash
# 连接到数据库
psql -U nexus_user -h localhost -d edit_business_db

# 执行迁移脚本
\i /path/to/001_add_user_fields.sql
```

### 3. 验证数据

```sql
-- 检查迁移前后数据量
SELECT
  '总用户数' as metric,
  COUNT(*) as count
FROM users

UNION ALL

SELECT
  '有 union_id 的用户数',
  COUNT(*)
FROM users
WHERE union_id IS NOT NULL;

-- 查看示例数据
SELECT
  id,
  auth_center_user_id,
  union_id,
  nickname,
  avatar_url,
  phone_number,
  email
FROM users
LIMIT 5;
```

### 4. 检查索引是否创建

```sql
-- 查看表的索引
\d users

-- 或使用系统表查询
SELECT
  indexname,
  indexdef
FROM pg_indexes
WHERE tablename = 'users';
```

### 5. 更新应用代码

迁移脚本执行完成后，重新编译并部署应用：

```bash
cd /Users/xia/Documents/GitHub/edit-business
go build -o edit-api cmd/server/main.go
```

## 回滚方案

如果迁移失败，可以使用以下回滚脚本：

```sql
-- 删除新字段（会同时删除索引）
ALTER TABLE users
  DROP COLUMN IF EXISTS union_id,
  DROP COLUMN IF EXISTS nickname,
  DROP COLUMN IF EXISTS avatar_url,
  DROP COLUMN IF EXISTS phone_number,
  DROP COLUMN IF EXISTS email;
```

## 性能对比

### 迁移前（JSONB）

```sql
-- 查询昵称（需要 JSON 解析）
SELECT * FROM users WHERE profile->>'nickname' = '张三';

-- 性能：较慢（需要解析 JSON）
```

### 迁移后（独立字段）

```sql
-- 查询昵称（直接查询字段）
SELECT * FROM users WHERE nickname = '张三';

-- 性能：快（可以使用 B-tree 索引）
```

## 代码改动

### Go 模型

```go
type User struct {
    // ... 原有字段

    // 新增字段
    UnionID     *string `gorm:"column:union_id;type:varchar(255);uniqueIndex"`
    Nickname    *string `gorm:"column:nickname;type:varchar(100)"`
    AvatarURL   *string `gorm:"column:avatar_url;type:varchar(500)"`
    PhoneNumber *string `gorm:"column:phone_number;type:varchar(255);uniqueIndex"`
    Email       *string `gorm:"column:email;type:varchar(255);uniqueIndex"`

    // Profile 保留用于低频字段
    Profile JSONB `gorm:"column:profile;type:jsonb;default:'{}'::jsonb"`
}
```

## 注意事项

1. **备份数据库** - 执行迁移前务必备份
2. **停机时间** - 迁移期间可能需要短暂停机
3. **测试环境** - 建议先在测试环境验证
4. **数据验证** - 迁移后务必验证数据完整性

## 联系人

如有问题，请联系开发团队。
