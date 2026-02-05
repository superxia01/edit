-- Drop index
DROP INDEX IF EXISTS idx_notes_user_id;

-- Remove foreign key constraint
ALTER TABLE notes DROP CONSTRAINT IF EXISTS fk_notes_user;

-- Remove user_id column
ALTER TABLE notes DROP COLUMN IF EXISTS user_id;
