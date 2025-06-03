package handlers

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
)

func createCharacterResponseMessage(char characters.Character, chatID int64) tgbotapi.Chattable {
	ratingStr          := fmt.Sprintf("%.1f", char.Rating)
	escapedName        := EscapeMarkdownV2(char.Name)
	escapedTitle       := EscapeMarkdownV2(char.Title)
	escapedDescription := EscapeMarkdownV2(char.Description)
	escapedRating      := EscapeMarkdownV2(ratingStr)
	
	caption := fmt.Sprintf(
		"*Имя:* %s\n"+
			"*ID:* `%d`\n"+
			"*Тайтл:* %s\n\n"+
			"*Описание:*\n%s\n\n"+
			"*Рейтинг:* %s",
		escapedName,
		char.ID,
		escapedTitle,
		escapedDescription,
		escapedRating,
	)
	
	if char.ImageURL != "" {
		photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(char.ImageURL))
		photoMsg.Caption = caption
		photoMsg.ParseMode = tgbotapi.ModeMarkdownV2
		return photoMsg
	}

	msg := tgbotapi.NewMessage(chatID, caption)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	return msg
}

func EscapeMarkdownV2(text string) string {
	escapeChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	var result strings.Builder
	
	for _, r := range text {
		char := string(r)
		found := false
		for _, ec := range escapeChars {
			if char == ec {
				result.WriteString("\\")
				result.WriteString(char)
				found = true
				break
			}
		}
		if !found {
			result.WriteString(char)
		}
	}
	return result.String()
}