ALTER TABLE characters DROP COLUMN deleted_at;
DROP INDEX IF EXISTS idx_characters_deleted_at;