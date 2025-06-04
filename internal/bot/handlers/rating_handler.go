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

// Структура для возврата результатов из HandleCallbackQuery
type CallbackResult struct {
	CallbackQueryAnswer *tgbotapi.CallbackConfig
	EditMessageText     *tgbotapi.EditMessageTextConfig
	EditMessageCaption  *tgbotapi.EditMessageCaptionConfig
	EditMessageMarkup   *tgbotapi.EditMessageReplyMarkupConfig
	NewMessage           tgbotapi.Chattable
}


// Обрабатывает команду /rate <id_or_name> <rating>
func HandleRateCommand(
	message    *tgbotapi.Message,
	charRepo   *characters.CharacterRepository,
	ratingRepo *ratings.Repository,
	) tgbotapi.Chattable {
		
	args  := strings.TrimSpace(message.CommandArguments())
	parts := strings.Fields(args)

	if len(parts) < 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Использование: /rate <ID или Имя> <оценка 1-10>\nНапример: `/rate 1 10`")
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		return msg
	}

	ratingStr  := parts[len(parts)-1]
	identifier := strings.Join(parts[:len(parts)-1], " ")

	ratingValue, err := strconv.Atoi(ratingStr)
	if err != nil || ratingValue < 1 || ratingValue > 10 {
		return tgbotapi.NewMessage(message.Chat.ID, "Оценка должна быть числом от 1 до 10.")
	}

	var char characters.Character
	var found bool
	charIDUint64, err := strconv.ParseUint(identifier, 10, 32)
	if err == nil {
		charID := uint(charIDUint64)
		char, found = charRepo.GetCharacterByID(charID)
	} else {
		char, found = charRepo.GetCharacterByNameOrAlt(identifier)
	}

	if !found {
		return tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Персонаж '%s' не найден.", EscapeMarkdownV2(identifier)))
	}

	// Сохраняем/обновляем оценку
	_, err = ratingRepo.SaveOrUpdateRating(message.From.ID, char.ID, ratingValue)
	if err != nil {
		log.Printf("Error saving rating via /rate command for char %d by user %d: %v", char.ID, message.From.ID, err)
		return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при сохранении вашей оценки.")
	}

	// Пересчитываем и обновляем средний рейтинг персонажа
	avgRating, count, err := ratingRepo.GetAverageRatingForCharacter(char.ID)
	if err != nil {
		log.Printf("Error recalculating average rating for char %d after /rate: %v", char.ID, err)
	} else {
		err = charRepo.UpdateCharacterRatingStats(char.ID, avgRating, count)
		if err != nil {
			log.Printf("Error updating character %d aggregate rating after /rate: %v", char.ID, err)
		}
	}

	return tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Спасибо! Вы оценили персонажа *%s* на *%d⭐*", EscapeMarkdownV2(char.Name), ratingValue))
}


// Обрабатывает нажатия на inline кнопки
func HandleCallbackQuery(
	api           *tgbotapi.BotAPI,
	callbackQuery *tgbotapi.CallbackQuery,
	charRepo      *characters.CharacterRepository,
	ratingRepo    *ratings.Repository,
) CallbackResult {
	
	var result CallbackResult
	userID    := callbackQuery.From.ID
	chatID    := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID
	data      := callbackQuery.Data

	answerCallback := tgbotapi.NewCallback(callbackQuery.ID, "")

	if strings.HasPrefix(data, CallbackPrefixRateAction) {
		charIDStr := strings.TrimPrefix(data, CallbackPrefixRateAction)
		charIDUint64, err := strconv.ParseUint(charIDStr, 10, 32)
		if err != nil {
			log.Printf("Error parsing charID from RateAction callback %s: %v", data, err)
			answerCallback.Text = "Ошибка: неверный ID персонажа."
			result.CallbackQueryAnswer = &answerCallback
			return result
		}
		characterID := uint(charIDUint64)

		var rows     [][]tgbotapi.InlineKeyboardButton
		var currentRow []tgbotapi.InlineKeyboardButton
		for i := 1; i <= 10; i++ {
			button := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%d⭐", i),
				fmt.Sprintf("%s%d_%d", CallbackPrefixSubmitRating, characterID, i),
			)
			currentRow = append(currentRow, button)
			if len(currentRow) == 5 || i == 10 { // 5 кнопок в ряду
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(currentRow...))
				currentRow = []tgbotapi.InlineKeyboardButton{}
			}
		}
		
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_rating_"+charIDStr)))

		ratingKeyboard          := tgbotapi.NewInlineKeyboardMarkup(rows...)
		editMarkup              := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, ratingKeyboard)
		result.EditMessageMarkup = &editMarkup
		answerCallback.Text      = "Выберите оценку:"
		
	} else if strings.HasPrefix(data, CallbackPrefixSubmitRating) {
		partsStr := strings.TrimPrefix(data, CallbackPrefixSubmitRating)
		parts    := strings.Split(partsStr, "_")
		if len(parts) != 2 {
			log.Printf("Error parsing SubmitRating callback data: %s", data)
			answerCallback.Text = "Ошибка: неверные данные для оценки."
			result.CallbackQueryAnswer = &answerCallback
			return result
		}

		charIDUint64, errChar  := strconv.ParseUint(parts[0], 10, 32)
		ratingValue, errRating := strconv.Atoi(parts[1])

		if errChar != nil || errRating != nil || ratingValue < 1 || ratingValue > 10 {
			log.Printf("Error parsing charID or ratingValue from SubmitRating callback %s: %v, %v", data, errChar, errRating)
			answerCallback.Text = "Ошибка: неверные данные оценки."
			result.CallbackQueryAnswer = &answerCallback
			return result
		}
		characterID := uint(charIDUint64)

		// Сохраняем/обновляем оценку
		_, err := ratingRepo.SaveOrUpdateRating(userID, characterID, ratingValue)
		if err != nil {
			log.Printf("Error saving rating from callback for char %d by user %d: %v", characterID, userID, err)
			answerCallback.Text = "Произошла ошибка при сохранении оценки."
			result.CallbackQueryAnswer = &answerCallback
			return result
		}

		// Пересчитываем и обновляем средний рейтинг персонажа
		avgRating, count, err := ratingRepo.GetAverageRatingForCharacter(characterID)
		charToUpdate, found   := charRepo.GetCharacterByID(characterID)

		if err != nil {
			log.Printf("Error recalculating average rating for char %d after callback: %v", characterID, err)
		} else if found {
			err = charRepo.UpdateCharacterRatingStats(characterID, avgRating, count)
			if err != nil {
				log.Printf("Error updating character %d aggregate rating after callback: %v", characterID, err)
			}
			charToUpdate.AverageRating = avgRating // Обновляем локальную копию для отображения
			charToUpdate.RatingCount   = count
		}

		answerCallback.Text = fmt.Sprintf("Вы оценили на %d⭐!", ratingValue)

		// Обновляем исходное сообщение (карточку персонажа) с новым рейтингом и кнопкой "Оценить"
		if found {
			updatedCardChattable := CreateChatCharacterResponseMessage(api, charToUpdate, chatID, userID, ratingRepo)

			if photoConfig, ok := updatedCardChattable.(tgbotapi.PhotoConfig); ok {
				editCap          := tgbotapi.NewEditMessageCaption(chatID, messageID, photoConfig.Caption)
				editCap.ParseMode = photoConfig.ParseMode
				if photoConfig.ReplyMarkup != nil {
					editCap.ReplyMarkup = photoConfig.ReplyMarkup.(*tgbotapi.InlineKeyboardMarkup)
				}
				result.EditMessageCaption = &editCap
			} else if messageConfig, ok := updatedCardChattable.(tgbotapi.MessageConfig); ok {
				editText          := tgbotapi.NewEditMessageText(chatID, messageID, messageConfig.Text)
				editText.ParseMode = messageConfig.ParseMode
				if messageConfig.ReplyMarkup != nil {
					editText.ReplyMarkup = messageConfig.ReplyMarkup.(*tgbotapi.InlineKeyboardMarkup)
				}
				result.EditMessageText = &editText
			}
		} else {
			// Если персонаж не найден (маловероятно здесь), просто убираем кнопки оценок
			editMarkup := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}})
			result.EditMessageMarkup = &editMarkup
		}
	} else {
		answerCallback.Text = "Неизвестное действие."
	}

	result.CallbackQueryAnswer = &answerCallback
	return result
}