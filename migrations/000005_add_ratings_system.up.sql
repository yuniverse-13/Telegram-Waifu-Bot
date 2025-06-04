CREATE TABLE user_ratings (
  id SERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  character_id INTEGER NOT NULL,
  rating INTEGER NOT NULL CHECK (rating >= 0 and rating <= 10),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  
  CONSTRAINT fk_character
    FOREIGN KEY(character_id)
    REFERENCES characters(id)
    ON DELETE CASCADE,
    
  UNIQUE (user_id, character_id)
);

CREATE INDEX idx_user_ratings_user_id ON user_ratings(user_id);
CREATE INDEX idx_user_ratings_character_id ON user_ratings(character_id);

ALTER TABLE characters DROP COLUMN IF EXISTS rating;
ALTER TABLE characters ADD COLUMN average_rating REAL DEFAULT 0.0;
ALTER TABLE characters ADD COLUMN rating_count INTEGER DEFAULT 0;