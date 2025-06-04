ALTER TABLE characters ADD COLUMN deleted_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS idx_characters_deleted_at ON characters(deleted_at);