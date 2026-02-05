-- Drop index
DROP INDEX IF EXISTS idx_bloggers_user_id;

-- Remove foreign key constraint
ALTER TABLE bloggers DROP CONSTRAINT IF EXISTS fk_bloggers_user;

-- Remove user_id column
ALTER TABLE bloggers DROP COLUMN IF EXISTS user_id;
