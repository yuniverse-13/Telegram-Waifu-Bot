package handlers

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
)

func HandleCharacterCommand(charRepo *characters.CharacterRepository, message *tgbotapi.Message) tgbotapi.Chattable {
	args := strings.TrimSpace(message.CommandArguments())
	var char characters.Character
	var found bool
		
	if args == "" {
		msg := "Пожалуйста, укажите имя или ID персонажа. Например: /character Фрирен или /character 1"
		return tgbotapi.NewMessage(message.Chat.ID, msg)
	}

	charID, err := strconv.Atoi(args)
	if err == nil && charID > 0 {
		char, found = charRepo.GetCharacterByID(charID)
	} else {
		char, found = charRepo.GetCharacterByNameOrAlt(args)
	}
	
	if !found {
		msg := fmt.Sprintf("Персонаж '%s' не найден.", args)
		return tgbotapi.NewMessage(message.Chat.ID, msg)
	}
	
	return createCharacterResponseMessage(char, message.Chat.ID)
}