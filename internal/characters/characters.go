package characters

import (
	"errors"
	"log"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Character struct {
	ID          int
	Name        string         `gorm:"uniqueIndex"`
	Title       string
	AltNames    pq.StringArray `gorm:"type:text[]"`
	Description string
	ImageURL    string         `gorm:"column:image_url"`
	Rating      float32
}

type CharacterRepository struct {
	db *gorm.DB
}

func NewCharacterRepository(db *gorm.DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

func (r *CharacterRepository) GetCharacterByID(id uint) (Character, bool) {
	var char Character
	
	result := r.db.First(&char, id)
		
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("GetCharacterByID: no character found for ID %d", id)
			return Character{}, false
		}
		log.Printf("GetCharacterByID: error getting character by ID %d: %v", id, result.Error)
		return Character{}, false
	}
	log.Printf("GetCharacterByID: Found character: %+v", char)
	return char, true
}

func (r *CharacterRepository) GetCharacterByNameOrAlt(name string) (Character, bool) {
	var char Character
	searchLower := strings.ToLower(name)
	
	result := r.db.Where("LOWER(name) = ?", searchLower).
		Or("? = ANY(LOWER(alt_names::text)::text[])", searchLower).First(&char)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("GetCharacterByNameOrAlt: no character found for name/alt '%s'", name)
			return Character{}, false
		}
		log.Printf("GetCharacterByNameOrAlt: error scanning row for name/alt '%s': %v", name, result.Error)
		return Character{}, false
	}
	log.Printf("GetCharacterByNameOrAlt: Found character: %+v for input '%s'", char, name)
	return char, true
}

func (r *CharacterRepository) GetRandomCharacter() (Character, bool) {
	var char Character
	
	result := r.db.Order("RANDOM()").Limit(1).First(&char)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("GetRandomCharacter: no characters in database to pick from.")
			return Character{}, false
		}
		log.Printf("GetRandomCharacter: error getting random character: %v", result.Error)
		return Character{}, false
	}
	log.Printf("GetRandomCharacter: Found random character: %+v", char)
	return char, true
}