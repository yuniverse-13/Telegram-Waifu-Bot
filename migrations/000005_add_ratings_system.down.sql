ALTER TABLE characters DROP COLUMN rating_count;
ALTER TABLE characters DROP COLUMN average_rating;
ALTER TABLE characters ADD COLUMN rating REAL;

DROP TABLE user_ratings;