-- Rollback Edit Business Database Schema

-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_bloggers_updated_at ON bloggers;
DROP TRIGGER IF EXISTS update_notes_updated_at ON notes;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_notes_author;
DROP INDEX IF EXISTS idx_notes_publish_date;
DROP INDEX IF EXISTS idx_notes_tags;
DROP INDEX IF EXISTS idx_notes_likes;
DROP INDEX IF EXISTS idx_bloggers_xhs_id;
DROP INDEX IF EXISTS idx_bloggers_followers;

-- Drop tables
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS bloggers;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "pgcrypto";
