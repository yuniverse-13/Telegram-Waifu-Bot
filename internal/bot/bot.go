package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/bot/handlers"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/ratings"
)

const (
	CommandStart           = "start"
	CommandCharacter       = "character"
	CommandRandomCharacter = "randomcharacter"
	CommandInfo            = "info"
	CommandRate            = "rate"
)

type Bot struct {
	api        *tgbotapi.BotAPI
	charRepo   *characters.CharacterRepository
	ratingRepo *ratings.Repository
}

func NewBot(token string, charRepo *characters.CharacterRepository, ratingRepo *ratings.Repository) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.Debug = true
	log.Printf("Authorized on account %s (@%s)", api.Self.FirstName, api.Self.UserName)

	return &Bot{
		api:        api,
		charRepo:   charRepo,
		ratingRepo: ratingRepo,
	}, nil
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)
	log.Println("Bot started listening for updates...")

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s (%d) in chat %d] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text)
			
			if update.Message.IsCommand() {
				response := b.handleCommand(update.Message)
				if response != nil {
					if _, err := b.api.Send(response); err != nil {
						log.Printf("Error sending message: %v. Response: %+v", err, response)
					}
				}
			} else {
				reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Я понимаю только команды.")
				b.api.Send(reply)
			}
		} else if update.CallbackQuery != nil {
			log.Printf("CallbackQuery: From %s (%d), Data: %s, MessageID: %d", update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID, update.CallbackQuery.Data, update.CallbackQuery.Message.MessageID)
			
			callbackResponse := handlers.HandleCallbackQuery(
				b.api,
				update.CallbackQuery,
				b.charRepo,
				b.ratingRepo,
			)
			
			if callbackResponse.CallbackQueryAnswer != nil {
				if _, err := b.api.Request(callbackResponse.CallbackQueryAnswer); err != nil {
					log.Printf("Error sending callback answer: %v", err)
				}
			}
			if callbackResponse.EditMessageMarkup != nil {
				if _, err := b.api.Request(callbackResponse.EditMessageMarkup); err != nil {
					log.Printf("Error editing message markup: %v", err)
				}
			}
			if callbackResponse.EditMessageText != nil {
				if _, err := b.api.Request(callbackResponse.EditMessageText); err != nil {
					log.Printf("Error editing message text: %v", err)
				}
			}
			if callbackResponse.EditMessageCaption != nil {
				if _, err := b.api.Request(callbackResponse.EditMessageCaption); err != nil {
					log.Printf("Error editing message caption: %v", err)
				}
			}
			if callbackResponse.NewMessage != nil {
				if _, err := b.api.Send(callbackResponse.NewMessage); err != nil {
					log.Printf("Error sending new message from callback: %v", err)
				}
			}
		}
	}
	log.Println("Update channel closed. Bot shutting down.")
	return nil
}

func (b *Bot) handleCommand(message *tgbotapi.Message) tgbotapi.Chattable {	
	var response tgbotapi.Chattable
	
	log.Printf("[%s (%d) in chat %d] Command: %s, Args: %s", message.From.UserName, message.From.ID, message.Chat.ID, message.Command(), message.CommandArguments())
	
	switch message.Command() {
	case CommandStart:
		response = handlers.HandleStartCommand(message)
	case CommandInfo:
		response = handlers.HandleInfoCommand(message)
	case CommandCharacter:
		response = handlers.HandleCharacterCommand(b.api, message, b.charRepo, b.ratingRepo)
	case CommandRandomCharacter:
		response = handlers.HandleRandomCharacterCommand(b.api, message, b.charRepo, b.ratingRepo)
	case CommandRate:
		response = handlers.HandleRateCommand(message, b.charRepo, b.ratingRepo)
	default:
		response = tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды.")
	}
	return response
}