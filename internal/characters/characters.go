package characters

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
)

type Character struct {
	ID          int
	Name        string
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
	err := r.db.QueryRow("SELECT id, name, alt_names, description, image_url, rating FROM characters WHERE id = $1", id).
		Scan(&char.ID, &char.Name, pq.Array(&char.AltNames), &char.Description, &char.ImageURL, &char.Rating)
		
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
		SELECT DISTINCT c.id, c.name, c.alt_names, c.description, c.image_url, c.rating
		FROM characters c
		LEFT JOIN unnest(COALESCE(c.alt_names, '{}'::text[])) AS alt_name_element ON true
		WHERE LOWER(c.name) = $1 OR LOWER(alt_name_element) = $1
		LIMIT 1;
	`
	logQuery := strings.ReplaceAll(strings.ReplaceAll(query, "\n", " "), "\t", "")
	log.Printf("Executing query: %s with param: %s", logQuery, searchLower)
	
	err := r.db.QueryRow(query, searchLower).
		Scan(&char.ID, &char.Name, pq.Array(&char.AltNames), &char.Description, &char.ImageURL, &char.Rating)
		
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
		SELECT id, name, alt_names, description, image_url, rating
		FROM characters
		ORDER BY RANDOM()
		LIMIT 1
	`
	
	err := r.db.QueryRow(query).
		Scan(&char.ID, &char.Name, pq.Array(&char.AltNames), &char.Description, &char.ImageURL, &char.Rating)
	
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

func (r *CharacterRepository) AddCharacter(char Character) (int, error) {
	var id int
	altNames := char.AltNames
	if altNames == nil {
		altNames = []string{}
	}
	
	query := "INSERT INTO characters (name, alt_names, description, image_url, rating) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	
	err := r.db.QueryRow(query, char.Name, pq.Array(char.AltNames), char.Description, char.ImageURL, char.Rating,).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("could not add character: %w", err)
	}
	return id, nil
}

// SeedInitialData - для заполнения БД начальными данными
/* func SeedInitialData(repo *CharacterRepository) {
	initialsCharacters := []Character{
		{
			ID: 1, Name: "Фрирен", AltNames: []string{"Frieren"},
			Description: "Эльфийская волшебница, бывший член отряда героев, победившего Короля Демонов.",
			ImageURL:    "https://i.pinimg.com/736x/87/b3/08/87b308eb264e7ee8c8bdc77774b87044.jpg", Rating: 9.8,
		},
		{
			ID: 2, Name: "Мегумин", AltNames: []string{"Мегу", "Megumin", "Megu"},
			Description: "Архимаг из клана Алых Демонов в фэнтезийном мире и первый человек, присоединившийся к группе Казумы.",
			ImageURL:    "https://i.pinimg.com/736x/be/f3/7a/bef37a878ab694f6722533b3c3222f1a.jpg", Rating: 9.9,
		},
		{
			ID: 3, Name: "Ферн", AltNames: []string{"Fern"},
			Description: "Юная, но удивительно серьёзная и ответственная волшебница, ученица Фрирен, которая часто выступает её своеобразным опекуном и голосом разума.",
			ImageURL:    "https://i.pinimg.com/736x/59/7b/2b/597b2b1f7906d73d87432150f181cb91.jpg", Rating: 9.5,
		},
		{
			ID: 4, Name: "Сильфиетта", AltNames: []string{"Сильфи", "Фиттс", "Sylphiette", "Sylphy"},
			Description: "Преданная эльфийка-полукровка, подруга детства и жена Рудеуса, выросшая из робкой девочки в талантливую волшебницу, также известная как Фиттс.",
			ImageURL:    "https://i.pinimg.com/736x/57/a7/4b/57a74bd586a877e4a75ad66e0b92fc5f.jpg", Rating: 8.8,
		},
		{
			ID: 5, Name: "Рокси", AltNames: []string{"Roxy"},
			Description: "Это талантливая и добрая демоница-мигурд, первая и глубоко уважаемая учительница магии Рудеуса.",
			ImageURL:    "https://i.pinimg.com/736x/5d/5c/93/5d5c93e86d2713c9dc7d758676382bc2.jpg", Rating: 8.9,
		},
	}
	log.Println("Starting initial data seeding...")
	
	for _, char := range initialsCharacters {
		_, found := repo.GetCharacterByNameOrAlt(char.Name)
		if !found {
			addedID, err := repo.AddCharacter(char)
			if err != nil {
				log.Printf("Error seeding character %s: %v", char.Name, err)
			} else {
				log.Printf("Seeded character: %s with ID: %d", char.Name, addedID)
			}
		} else {
			log.Printf("Character %s already exists or found by alt name, skipping seed.", char.Name)
		}
	}
	log.Println("Initial data seeding complete.")
} */