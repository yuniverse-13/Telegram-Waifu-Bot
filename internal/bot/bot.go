package bot

import (
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
)

type Bot struct {
	api *tgbotapi.BotAPI
	// Здесь можно добавить другие зависимости, например, сервис для работы с БД персонажей
}

func NewBot(token string) (*Bot, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	api.Debug = true
	log.Printf("Authorized on account %s (@%s)", api.Self.FirstName, api.Self.UserName)

	return &Bot{api: api}, nil
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
			b.handleCommand(update.Message)
		} else {
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Я понимаю только команды.")
			b.api.Send(reply)
		}
	}

	log.Println("Update channel closed. Bot shutting down.")
	return nil
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	var response tgbotapi.Chattable

	switch message.Command() {
	case "start":
		responseText := fmt.Sprintf("Привет, %s! Я твой Waifu бот. Используй /character, чтобы увидеть персонажа.", message.From.FirstName)
		response = tgbotapi.NewMessage(message.Chat.ID, responseText)
	case "character":
		char := characters.GetSampleCharacter()

		caption := fmt.Sprintf("Имя: %s\nОписание: %s\nРейтинг: %d/10", char.Name, char.Description, char.Rating)

		if char.ImageURL != "" {
			photoMsg := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL(char.ImageURL))
			photoMsg.Caption = caption
			response = photoMsg
		} else {
			response = tgbotapi.NewMessage(message.Chat.ID, caption)
		}

	case "megumin":
		char := characters.GetAnotherSampleCharacter()
		caption := fmt.Sprintf("Имя: %s\nОписание: %s\nРейтинг: %d/10", char.Name, char.Description, char.Rating)
		if char.ImageURL != "" {
			photoMsg := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL(char.ImageURL))
			photoMsg.Caption = caption
			response = photoMsg
		} else {
			response = tgbotapi.NewMessage(message.Chat.ID, caption)
		}

	default:
		response = tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такую команду.")
	}

	if response != nil {
		if _, err := b.api.Send(response); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}
