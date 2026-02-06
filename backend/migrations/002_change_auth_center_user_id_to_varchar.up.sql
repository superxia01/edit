-- =====================================================
-- Edit Business: 用户表优化 - auth_center_user_id 改为 VARCHAR
-- 版本: 2026-02-06
-- 说明: 与 PR、Pixel、Quote 系统保持一致，使用 string 类型存储 UUID
-- =====================================================

-- 1. 修改 auth_center_user_id 列类型从 UUID 改为 VARCHAR(255)
ALTER TABLE users
ALTER COLUMN auth_center_user_id TYPE VARCHAR(255);

-- 2. 删除旧的 UUID 类型约束（如果有）
-- ALTER TABLE users
-- DROP CONSTRAINT IF EXISTS users_auth_center_user_id_key;

-- 3. 重新添加唯一索引（使用 VARCHAR 类型）
DROP INDEX IF EXISTS idx_users_auth_center_user_id;
CREATE UNIQUE INDEX idx_users_auth_center_user_id ON users(auth_center_user_id);

-- 4. 修改 id 列类型从 UUID 改为 VARCHAR(255)（如果需要）
ALTER TABLE users
ALTER COLUMN id TYPE VARCHAR(255);

-- 5. 添加注释
COMMENT ON COLUMN users.auth_center_user_id IS '账号中心用户 ID（从 auth-center 同步，字符串格式）';
COMMENT ON COLUMN users.id IS '用户唯一标识（本地生成，格式：user-{timestamp}）';

-- 验证迁移结果
SELECT
    'Migration completed' as status,
    COUNT(*) as total_users,
    COUNT(auth_center_user_id) as users_with_auth_center_id
FROM users;
