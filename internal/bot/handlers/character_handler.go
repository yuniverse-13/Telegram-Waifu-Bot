package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/ratings"
)

func HandleCharacterCommand(
	api        *tgbotapi.BotAPI,
	message    *tgbotapi.Message,
	charRepo   *characters.CharacterRepository,
	ratingRepo *ratings.Repository,
	) tgbotapi.Chattable {
	
	args := strings.TrimSpace(message.CommandArguments())
	var char  characters.Character
	var found bool
		
	if args == "" {
		errorText    := "Пожалуйста, укажите ID или имя персонажа. Например: `/character 1` или `/character Фрирен`"
		msg          := tgbotapi.NewMessage(message.Chat.ID, errorText)
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		return msg
	}

	charIDUint64, err := strconv.ParseUint(args, 10, 32)
	if err == nil {
		charID := uint(charIDUint64)
		log.Printf("HandleCharacterCommand: searching by ID: %d", charID)
		char, found = charRepo.GetCharacterByID(charID)
	} else {
		log.Printf("HandleCharacterCommand: searching by Name/AltName: %s", args)
		char, found = charRepo.GetCharacterByNameOrAlt(args)
	}
	
	if !found {
		log.Printf("HandleCharacterCommand: Character not found for query: %s", args)
		msg := fmt.Sprintf("Персонаж '%s' не найден.", args)
		return tgbotapi.NewMessage(message.Chat.ID, msg)
	}
	
	log.Printf("HandleCharacterCommand: Character found: %s", char.Name)
	return CreateChatCharacterResponseMessage(api, char, message.Chat.ID, message.From.ID, ratingRepo)
}