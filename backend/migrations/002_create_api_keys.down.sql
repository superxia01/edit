-- Drop API Keys table
DROP INDEX IF EXISTS idx_api_keys_deleted_at;
DROP INDEX IF EXISTS idx_api_keys_key;
DROP INDEX IF EXISTS idx_api_keys_user_id;
DROP TABLE IF EXISTS api_keys;
