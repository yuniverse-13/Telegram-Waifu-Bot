package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
)

func createCharacterResponseMessage(char characters.Character, chatID int64) tgbotapi.Chattable {
	caption := fmt.Sprintf(
		"Имя: %s\nID: %d\n\nТайтл: %s\n\nОписание: %s\n\nРейтинг: %.1f",
		char.Name, char.ID, char.Title, char.Description, char.Rating,
	)
	
	if char.ImageURL != "" {
		photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(char.ImageURL))
		photoMsg.Caption = caption
		return photoMsg
	}
	msg := tgbotapi.NewMessage(chatID, "Изображение отсутствует\n\n" + caption)
	return msg
}