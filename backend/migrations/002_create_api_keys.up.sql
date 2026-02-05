-- Create API Keys table for plugin authentication
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_used TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index for faster user queries
CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_key ON api_keys(key) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_api_keys_deleted_at ON api_keys(deleted_at);

-- Add comment
COMMENT ON TABLE api_keys IS 'API keys for Chrome plugin authentication';
COMMENT ON COLUMN api_keys.key IS 'API key string (stored as plain text for plugin validation)';
COMMENT ON COLUMN api_keys.is_active IS 'Whether the API key is active';
COMMENT ON COLUMN api_keys.last_used IS 'Last timestamp when this key was used';
COMMENT ON COLUMN api_keys.expires_at IS 'Optional expiration date';
