package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
)

type Bot struct {
	api *tgbotapi.BotAPI
	// Здесь можно добавить другие зависимости, например, сервис для работы с БД персонажей
}

func NewBot(token string) (*Bot, error) {
	if token == "" {
		return nil, errors.New("NewBot: токен пустой")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать API бота: %w", err)
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
	var err error

	switch message.Command() {
	case "start":
		responseText := fmt.Sprintf("Привет, %s! Я твой Waifu бот. Используй /character <имя> или /randomcharacter, чтобы увидеть персонажа.", message.From.FirstName)
		response = tgbotapi.NewMessage(message.Chat.ID, responseText)
		
	case "character":
		charName := strings.TrimSpace(message.CommandArguments())
		if charName == "" {
			response = tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста, укажи имя персонажа после команды /character.\nНапример: /character Фрирен")
		} else {
			char, found := characters.GetCharacterByName(charName)
			if found {
				caption := fmt.Sprintf("Имя: %s\nОписание: %s\nРейтинг: %.1f / 10", char.Name, char.Description, char.Rating)
				if char.ImageURL != "" {
					photoMsg := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL(char.ImageURL))
					photoMsg.Caption = caption
					response = photoMsg
				} else {
					response = tgbotapi.NewMessage(message.Chat.ID, caption + "\n(Изображение отсутствует)")
				}
			} else {
				response = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Персонаж '%s' не найден.", charName))
			}
		}

	case "randomcharacter":
		char, found := characters.GetRandomCharacter()
		if found {
			caption := fmt.Sprintf("Случайный персонаж!\nИмя: %s\nОписание: %s\nРейтинг: %.1f / 10", char.Name, char.Description, char.Rating)
			if char.ImageURL != "" {
				photoMsg := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL(char.ImageURL))
				photoMsg.Caption = caption
				response = photoMsg
			} else {
				response = tgbotapi.NewMessage(message.Chat.ID, caption+"\n(Изображение отсутствует)")
			}
		} else {
			response = tgbotapi.NewMessage(message.Chat.ID, "В базе данных нет персонажей для отображения.")
		}
		
	default:
		response = tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды.")
	}
	
	if response != nil {
		_, err = b.api.Send(response)
		if err != nil {
			log.Printf("Ошибка отправки ответа: %v", err)
		}
	}
}
