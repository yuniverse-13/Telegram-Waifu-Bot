CREATE TABLE characters (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  alt_names TEXT[],
  description TEXT,
  image_url VARCHAR(512),
  rating REAL
);

CREATE INDEX idx_characters_name ON characters (name);
CREATE INDEX idx_characters_alt_names ON characters USING GIN (alt_names);