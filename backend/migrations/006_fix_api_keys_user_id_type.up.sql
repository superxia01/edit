-- =====================================================
-- Fix api_keys.user_id type to match users.id (VARCHAR)
-- 新用户创建 API Key 失败原因：users.id 已改为 VARCHAR，api_keys.user_id 仍为 UUID
-- =====================================================

-- 1. 删除 api_keys 到 users 的外键约束（名称可能因环境而异）
DO $$
DECLARE r RECORD;
BEGIN
  FOR r IN (SELECT conname FROM pg_constraint
            WHERE conrelid = 'api_keys'::regclass AND confrelid = 'users'::regclass)
  LOOP
    EXECUTE 'ALTER TABLE api_keys DROP CONSTRAINT ' || quote_ident(r.conname);
  END LOOP;
END $$;

-- 2. 修改 user_id 类型为 VARCHAR(255)
ALTER TABLE api_keys ALTER COLUMN user_id TYPE VARCHAR(255) USING user_id::TEXT;

-- 3. 重新添加外键约束
ALTER TABLE api_keys ADD CONSTRAINT api_keys_user_id_fkey
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
