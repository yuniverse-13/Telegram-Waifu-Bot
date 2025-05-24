package characters

import (
	"database/sql"
	"log"
	"strings"

	"github.com/lib/pq"
)

type Character struct {
	ID          int
	Name        string
	Title       string
	AltNames    []string
	Description string
	ImageURL    string
	Rating      float32
}

type CharacterRepository struct {
	db *sql.DB
}

func NewCharacterRepository(db *sql.DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

func (r *CharacterRepository) GetCharacterByID(id int) (Character, bool) {
	char := Character{}
	err := r.db.QueryRow("SELECT id, name, title, alt_names, description, image_url, rating FROM characters WHERE id = $1", id).
		Scan(&char.ID, &char.Name, &char.Title, pq.Array(&char.AltNames), &char.Description, &char.ImageURL, &char.Rating)
		
	if err != nil {
		if err == sql.ErrNoRows {
			return Character{}, false
		}
		log.Printf("Error getting character by ID %d: %v", id, err)
		return Character{}, false
	}
	return char, true
}

func (r *CharacterRepository) GetCharacterByNameOrAlt(name string) (Character, bool) {
	char := Character{}
	searchLower := strings.ToLower(name)
	
	log.Printf("GetCharacterByNameOrAlt: input name='%s', searchLower='%s'", name, searchLower)
	
	query := `
		SELECT DISTINCT c.id, c.name, c.title, c.alt_names, c.description, c.image_url, c.rating
		FROM characters c
		LEFT JOIN unnest(COALESCE(c.alt_names, '{}'::text[])) AS alt_name_element ON true
		WHERE LOWER(c.name) = $1 OR LOWER(alt_name_element) = $1
		LIMIT 1;
	`
	logQuery := strings.ReplaceAll(strings.ReplaceAll(query, "\n", " "), "\t", "")
	log.Printf("Executing query: %s with param: %s", logQuery, searchLower)
	
	err := r.db.QueryRow(query, searchLower).
		Scan(&char.ID, &char.Name, &char.Title, pq.Array(&char.AltNames), &char.Description, &char.ImageURL, &char.Rating)
		
	if err != nil	{
		if err == sql.ErrNoRows	{
			log.Printf("GetCharacterByNameOrAlt: No rows found for '%s'", searchLower)
			return Character{}, false
		}
		log.Printf("GetCharacterByNameOrAlt: Error scanning row for '%s': %v", searchLower, err)
		return Character{}, false
	}
	log.Printf("GetCharacterByNameOrAlt: Found character: %+v", char)
	return char, true
}

func (r *CharacterRepository) GetRandomCharacter() (Character, bool) {
	char := Character{}
	
	query := `
		SELECT id, name, title, alt_names, description, image_url, rating
		FROM characters
		ORDER BY RANDOM()
		LIMIT 1
	`
	
	err := r.db.QueryRow(query).
		Scan(&char.ID, &char.Name, &char.Title, pq.Array(&char.AltNames), &char.Description, &char.ImageURL, &char.Rating)
	
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No characters in database to pick random from.")
			return Character{}, false
		}
		log.Printf("Error getting random character: %v", err)
		return Character{}, false
	}
	return char, true
}