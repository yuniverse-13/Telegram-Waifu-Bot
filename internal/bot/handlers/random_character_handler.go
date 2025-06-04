package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/ratings"
)

func HandleRandomCharacterCommand(
	api        *tgbotapi.BotAPI,
	message    *tgbotapi.Message,
	charRepo   *characters.CharacterRepository,
	ratingRepo *ratings.Repository,
	) tgbotapi.Chattable {
		
	char, found := charRepo.GetRandomCharacter()
	if !found {
		msg := "Не удалось получить случайного персонажа. Возможно, в базе данных пока нет записей."
		return tgbotapi.NewMessage(message.Chat.ID, msg)
	}
	return CreateChatCharacterResponseMessage(api, char, message.Chat.ID, message.From.ID, ratingRepo)
}