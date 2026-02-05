-- Create user settings table for plugin control
CREATE TABLE IF NOT EXISTS user_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    collection_enabled BOOLEAN NOT NULL DEFAULT false,
    collection_daily_limit INTEGER NOT NULL DEFAULT 500,
    collection_batch_limit INTEGER NOT NULL DEFAULT 50,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Add comment
COMMENT ON TABLE user_settings IS 'User settings for plugin collection control';
COMMENT ON COLUMN user_settings.collection_enabled IS '允许采集开关（用户可修改）';
COMMENT ON COLUMN user_settings.collection_daily_limit IS '每天采集上限（只读，开发者配置）';
COMMENT ON COLUMN user_settings.collection_batch_limit IS '单次采集上限（只读，开发者配置）';
