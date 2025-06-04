package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	welcomeMessage  = "Привет, %s\\! Я твой Waifu бот\\. Используй `/info`, чтобы увидеть список команд\\."
	helpMessage     = "Доступные команды:\n" +
		"`/start` \\- Приветственное сообщение\n" +
		"/character \\<ID или имя\\> \\- Поиск персонажа по ID или имени\\. Например: `/character 1` или `/character Фрирен`\n" +
		"`/randomcharacter` \\- Показать случайного персонажа\n" +
		"/rate \\<ID или имя\\> \\<оценка 1\\-10\\> \\- Оценить персонажа\\. Например: `/rate 1 10` или `/rate Фрирен 10`\n" +
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