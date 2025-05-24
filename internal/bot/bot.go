package bot

import (
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/bot/handlers"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
)

const (
	CommandStart           = "start"
	CommandCharacter       = "character"
	CommandRandomCharacter = "randomcharacter"
	CommandInfo            = "info"
)

type Bot struct {
	api *tgbotapi.BotAPI
	charRepo *characters.CharacterRepository
}

func NewBot(token string, charRepo *characters.CharacterRepository) (*Bot, error) {
	if token == "" {
		return nil, errors.New("NewBot: токен пустой")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать API бота: %w", err)
	}

	api.Debug = true
	log.Printf("Authorized on account %s (@%s)", api.Self.FirstName, api.Self.UserName)

	return &Bot{api: api, charRepo: charRepo,}, nil
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	log.Println("Bot started listening for updates...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s (%d)] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

		if update.Message.IsCommand() {
			b.handlerCommand(update.Message)
		} else {
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Я понимаю только команды.")
			b.api.Send(reply)
		}
	}
	log.Println("Update channel closed. Bot shutting down.")
	return nil
}

func (b *Bot) handlerCommand(message *tgbotapi.Message) {
	var response tgbotapi.Chattable
	var err error
	
	switch message.Command() {
	case CommandStart:
		response = handlers.HandleStartCommand(message)
	case CommandInfo:
		response = handlers.HandleInfoCommand(message)
	case CommandCharacter:
		response = handlers.HandleCharacterCommand(b.charRepo, message)
	case CommandRandomCharacter:
		response = handlers.HandleRandomCharacterCommand(b.charRepo, message)
	default:
		response = tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды.")
	}
	
	if response != nil {
		if _, err = b.api.Send(response); err != nil {
			log.Printf("Ошибка отправки ответа: %s", err)
		}
	}
}