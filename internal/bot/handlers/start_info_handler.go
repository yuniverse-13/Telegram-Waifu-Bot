package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	welcomeMessage  = "Привет, %s\\! Я твой Waifu бот\\. Используй `/info`, чтобы увидеть список команд\\."
	helpMessage     = "Доступные команды:\n" +
		"`/start` \\- Приветственное сообщение\n" +
		"`/character \\<ID или имя\\>` \\- Поиск персонажа по ID или имени\\. Пример: `/character 1` или `/character Фрирен`\n" +
		"`/randomcharacter` \\- Показать случайного персонажа\n" +
		"`/info` \\- Показать это сообщение со списком команд"
)

func HandleStartCommand(message *tgbotapi.Message) tgbotapi.Chattable {
	text := fmt.Sprintf(welcomeMessage, message.From.FirstName)
	msg  := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	return msg
}

func HandleInfoCommand(message *tgbotapi.Message) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(message.Chat.ID, helpMessage)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	return msg
}