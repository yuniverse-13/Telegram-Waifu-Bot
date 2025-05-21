package bot

import (
	"errors"
	"fmt"
	"log"
	"strconv"
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
	
	args := strings.TrimSpace(message.CommandArguments())

	switch message.Command() {
	case "start":
		responseText := fmt.Sprintf("Привет, %s! Я твой Waifu бот. Используй /character <имя или ID> или /randomcharacter, чтобы увидеть персонажа.\n/info - Показать сообщение со списком команд", message.From.FirstName)
		response = tgbotapi.NewMessage(message.Chat.ID, responseText)
		
	case "character":
		if args == "" {
			responseText := fmt.Sprintln("Укажите имя или ID персонажа. Например: /character Фрирен или /character 1")
			response = tgbotapi.NewMessage(message.Chat.ID, responseText)
		}
		
		var char characters.Character
		var found bool
		
		charID, convErr := strconv.Atoi(args)
		if convErr == nil {
			char, found = characters.GetCharacterByID(charID)
		} else {
			char, found = characters.GetCharacterByNameOrAlt(args)
		}
		
		if found {
			caption := fmt.Sprintf("Имя: %s\nОписание: %s\nРейтинг: %.1f из 10", char.Name, char.Description, char.Rating)
			if char.ImageURL != "" {
				photoMsg := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL(char.ImageURL))
				photoMsg.Caption = caption
				response = photoMsg
			} else {
				response = tgbotapi.NewMessage(message.Chat.ID, caption+"\n(Изображение отсутствует)")
			}
		} else {
			if convErr == nil {
				response = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Персонаж с ID %d не найден.", charID))
			} else {
				response = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Персонаж '%s' не найден.", args))
			}
		}

	case "randomcharacter":
		char, found := characters.GetRandomCharacter()
		
		if found {
			caption := fmt.Sprintf("Случайный персонаж!\nИмя: %s\nОписание: %s\nРейтинг: %.1f из 10", char.Name, char.Description, char.Rating)
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
		
	case "info":
		infoText := `
Доступные команды:
/start - Приветственное сообщение
/character <ID или имя> - Поиск персонажа по ID или имени. Пример: /character 1 или /character Фрирен
/randomcharacter - Показать случайного персонажа
/info - Показать это сообщение со списком команд
`
		response = tgbotapi.NewMessage(message.Chat.ID, infoText)
		
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