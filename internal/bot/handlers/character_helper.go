package handlers

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/ratings"
)

const (
	CallbackPrefixRateAction = "rate_action_"
	CallbackPrefixSubmitRating  = "submit_rating_"
)

var escapeChars = []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

func EscapeMarkdownV2(text string) string {
	var result strings.Builder
	
	for _, r := range text {
		char := string(r)
		isEscapeChar := false
		for _, ec := range escapeChars {
			if char == ec {
				result.WriteString("\\")
				result.WriteString(char)
				isEscapeChar = true
				break
			}
		}
		if !isEscapeChar {
			result.WriteString(char)
		}
	}
	return result.String()
}


// Формирует сообщение с карточкой персонажа
func CreateChatCharacterResponseMessage(
	api        *tgbotapi.BotAPI,
	char        characters.Character,
	chatID      int64,
	userID      int64,
	ratingRepo *ratings.Repository,
) tgbotapi.Chattable {

	userSpecificRating, err := ratingRepo.GetUserRatingForCharacter(userID, char.ID)
	var userRatingStr string
	if err != nil {
		log.Printf("Error fetching user's (%d) rating for char %d: %v", userID, char.ID, err)
		userRatingStr = "Ошибка загрузки"
	} else if userSpecificRating != nil {
		userRatingStr = fmt.Sprintf("%d⭐", userSpecificRating.Rating)
	} else {
		userRatingStr = "Вы не оценили"
	}

	char.AverageRating, char.RatingCount, err = ratingRepo.GetAverageRatingForCharacter(char.ID)
	if err != nil {
		log.Printf("Error fetching average rating for char %d: %v", char.ID, err)
		char.AverageRating = 0.0
		char.RatingCount = 0
	}

	averageRatingStr := fmt.Sprintf("%.1f", char.AverageRating)

	caption := fmt.Sprintf(
		"*Имя:* `%s`\n"+
			"*Тайтл:* `%s`\n"+
			"*ID:* `%d`\n\n"+
			"*Рейтинг:* %s⭐ \\(голосов: %d\\)\n"+
			"*Ваша оценка:* %s\n\n"+
			"*Описание:*\n%s",
		EscapeMarkdownV2(char.Name),
		EscapeMarkdownV2(char.Title),
		char.ID,
		EscapeMarkdownV2(averageRatingStr),
		char.RatingCount,
		EscapeMarkdownV2(userRatingStr),
		EscapeMarkdownV2(char.Description),
	)
	
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Оценить ✨", fmt.Sprintf("%s%d", CallbackPrefixRateAction, char.ID)),
		),
	)

	var response tgbotapi.Chattable
	if char.ImageURL != "" {
		photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(char.ImageURL))
		photoMsg.Caption = caption
		photoMsg.ParseMode = tgbotapi.ModeMarkdownV2
		photoMsg.ReplyMarkup = &inlineKeyboard
		response = photoMsg
	} else {
		msg := tgbotapi.NewMessage(chatID, caption)
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		msg.ReplyMarkup = &inlineKeyboard
		response = msg
	}
	return response
}