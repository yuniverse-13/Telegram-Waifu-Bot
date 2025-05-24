package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
)

func HandleRandomCharacterCommand(charRepo *characters.CharacterRepository, message *tgbotapi.Message) tgbotapi.Chattable {
	char, found := charRepo.GetRandomCharacter()
	
	if !found {
		msg := "Не удалось получить случайного персонажа. Возможно, в базе данных пока нет записей."
		return tgbotapi.NewMessage(message.Chat.ID, msg)
	}
	return createCharacterResponseMessage(char, message.Chat.ID)
}