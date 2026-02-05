-- Add source column to notes table
-- Migration date: 2026-02-05
-- Description: Add source column to distinguish between single note capture and batch blogger capture

ALTER TABLE notes ADD COLUMN IF NOT EXISTS source VARCHAR(20) DEFAULT 'single';

-- Add comment
COMMENT ON COLUMN notes.source IS 'Capture source: single = single note capture, batch = blogger batch capture';

-- Update existing records (optional, set default for old records)
UPDATE notes SET source = 'single' WHERE source IS NULL;

-- Create index for filtering by source
CREATE INDEX IF NOT EXISTS idx_notes_source ON notes(source);
