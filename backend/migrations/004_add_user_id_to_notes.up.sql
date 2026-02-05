-- Add user_id to notes table
ALTER TABLE notes ADD COLUMN IF NOT EXISTS user_id UUID NOT NULL DEFAULT gen_random_uuid();
ALTER TABLE notes ADD CONSTRAINT fk_notes_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);
