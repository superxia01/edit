-- Add user_id to bloggers table
ALTER TABLE bloggers ADD COLUMN IF NOT EXISTS user_id UUID NOT NULL DEFAULT gen_random_uuid();
ALTER TABLE bloggers ADD CONSTRAINT fk_bloggers_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_bloggers_user_id ON bloggers(user_id);
