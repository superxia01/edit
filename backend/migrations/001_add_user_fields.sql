-- ============================================
-- 迁移脚本：将用户信息从 JSONB 迁移到独立字段
-- ============================================
-- 目的：提升查询性能，避免 JSONB 解析开销，可以建立唯一索引
-- 日期：2026-02-06
-- ============================================

-- ============================================
-- 第1步：添加新字段
-- ============================================
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS union_id VARCHAR(255) UNIQUE,
  ADD COLUMN IF NOT EXISTS nickname VARCHAR(100),
  ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500),
  ADD COLUMN IF NOT EXISTS phone_number VARCHAR(255) UNIQUE,
  ADD COLUMN IF NOT EXISTS email VARCHAR(255) UNIQUE;

-- ============================================
-- 第2步：从 JSONB 迁移数据到新字段
-- ============================================
-- 迁移 unionId
UPDATE users
SET union_id = (profile->>'unionId')::VARCHAR(255)
WHERE profile IS NOT NULL
  AND profile ? 'unionId'
  AND (profile->>'unionId') IS NOT NULL
  AND (profile->>'unionId') != '';

-- 迁移 nickname
UPDATE users
SET nickname = (profile->>'nickname')::VARCHAR(100)
WHERE profile IS NOT NULL
  AND profile ? 'nickname'
  AND (profile->>'nickname') IS NOT NULL
  AND (profile->>'nickname') != '';

-- 迁移 avatarUrl
UPDATE users
SET avatar_url = (profile->>'avatarUrl')::VARCHAR(500)
WHERE profile IS NOT NULL
  AND profile ? 'avatarUrl'
  AND (profile->>'avatarUrl') IS NOT NULL
  AND (profile->>'avatarUrl') != '';

-- 迁移 phoneNumber
UPDATE users
SET phone_number = (profile->>'phoneNumber')::VARCHAR(255)
WHERE profile IS NOT NULL
  AND profile ? 'phoneNumber'
  AND (profile->>'phoneNumber') IS NOT NULL
  AND (profile->>'phoneNumber') != '';

-- 迁移 email
UPDATE users
SET email = (profile->>'email')::VARCHAR(255)
WHERE profile IS NOT NULL
  AND profile ? 'email'
  AND (profile->>'email') IS NOT NULL
  AND (profile->>'email') != '';

-- ============================================
-- 第3步：创建索引
-- ============================================
-- union_id 索引（已由 UNIQUE 约束自动创建）
CREATE INDEX IF NOT EXISTS idx_users_nickname ON users(nickname) WHERE nickname IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number) WHERE phone_number IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE email IS NOT NULL;

-- ============================================
-- 第4步：数据验证
-- ============================================
-- 检查迁移前后的数据量
SELECT
  '总用户数' as metric,
  COUNT(*) as count
FROM users

UNION ALL

SELECT
  '有 union_id 的用户数',
  COUNT(*)
FROM users
WHERE union_id IS NOT NULL

UNION ALL

SELECT
  '有 nickname 的用户数',
  COUNT(*)
FROM users
WHERE nickname IS NOT NULL

UNION ALL

SELECT
  '有 avatar_url 的用户数',
  COUNT(*)
FROM users
WHERE avatar_url IS NOT NULL

UNION ALL

SELECT
  '有 phone_number 的用户数',
  COUNT(*)
FROM users
WHERE phone_number IS NOT NULL

UNION ALL

SELECT
  '有 email 的用户数',
  COUNT(*)
FROM users
WHERE email IS NOT NULL;

-- ============================================
-- 第5步：查看表结构（验证新字段已添加）
-- ============================================
\d users

-- ============================================
-- 回滚脚本（如果需要回滚）
-- ============================================
/*
-- 删除新字段（会同时删除索引）
ALTER TABLE users
  DROP COLUMN IF EXISTS union_id,
  DROP COLUMN IF EXISTS nickname,
  DROP COLUMN IF EXISTS avatar_url,
  DROP COLUMN IF EXISTS phone_number,
  DROP COLUMN IF EXISTS email;
*/
