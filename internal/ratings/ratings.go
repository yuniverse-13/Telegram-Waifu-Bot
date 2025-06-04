package ratings

import (
	"errors"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GORM модель для оценок
type UserRating struct {
	gorm.Model
	UserID      int64 `gorm:"uniqueIndex:idx_user_character_rating"`
	CharacterID uint  `gorm:"uniqueIndex:idx_user_character_rating"`
	Rating      int
}

// Repository для работы с оценками
type Repository struct {
	db *gorm.DB
}

// Конструктор для RatingRepository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Сохраняет или обновляет оценку пользователя для персонажа
func (r *Repository) SaveOrUpdateRating(userID int64, characterID uint, ratingValue int) (*UserRating, error) {
	rating := UserRating{
		UserID:      userID,
		CharacterID: characterID,
		Rating:      ratingValue,
	}
	
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "character_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"rating", "updated_at"}),
	}).Create(&rating).Error
	
	if err != nil {
		log.Printf("Error saving/updating rating for user %d, char %d: %v", userID, characterID, err)
		return nil, err
	}
	log.Printf("Saved/Updated rating: UserID %d, CharID %d, Rating %d", userID, characterID, ratingValue)
	return &rating, nil
}

// Вычисляет средний рейтинг и количество оценок для персонажа
func (r *Repository) GetAverageRatingForCharacter(characterID uint) (avgRating float32, count int, err error) {
	var result struct {
		Avg   float32
		Count int64
	}
	
	err = r.db.Model(&UserRating{}).
		Select("COALESCE(AVG(rating), 0.0) as avg, COUNT(rating) as count").
		Where("character_id = ? AND deleted_at IS NULL", characterID).
		Scan(&result).Error
		
	if err != nil {
		log.Printf("Error calculating average rating for char %d: %v", characterID, err)
		return 0.0, 0, err
	}
	return result.Avg, int(result.Count), nil
}

//Получает оценку конкретного пользователя для конкретного персонажа
func (r *Repository) GetUserRatingForCharacter(userID int64, characterID uint) (*UserRating, error) {
	var rating UserRating
	
	err := r.db.Where("user_id = ? AND character_id = ? AND deleted_at IS NULL", userID, characterID).First(&rating).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rating, nil
}