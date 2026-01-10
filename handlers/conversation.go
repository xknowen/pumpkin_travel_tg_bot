package handlers

import (
	"fmt"
	"pumpkin_travel_tg_bot/models"
	"pumpkin_travel_tg_bot/services"
	"pumpkin_travel_tg_bot/utils"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type ConversationHandler struct {
	commandHandler *CommandHandler
	formService    *services.FormService
}

func NewConversationHandler(
	cmdHandler *CommandHandler,
	formService *services.FormService,
) *ConversationHandler {
	return &ConversationHandler{
		commandHandler: cmdHandler,
		formService:    formService,
	}
}

func (ch *ConversationHandler) HandleMessage(update tgbotapi.Update) {
	userID := update.Message.From.ID
	logrus.WithFields(logrus.Fields{
		"user_id": userID,
		"text":    update.Message.Text,
		"step":    ch.commandHandler.userStep[userID],
	}).Info("Обработка сообщения")

	// Получаем состояние пользователя
	state, step, exists := ch.commandHandler.GetUserState(userID)
	if !exists {
		logrus.WithField("user_id", userID).Warn("Пользователь не в диалоге, показываем помощь")
		// Пользователь не в диалоге
		ch.commandHandler.HandleHelp(update)
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id": userID,
		"step":    step,
	}).Info("Обработка шага диалога")

	// Обрабатываем сообщение в зависимости от шага
	switch step {
	case 1:
		ch.handleDestinationType(update, state, userID)
	case 2:
		ch.handleCountries(update, state, userID)
	case 3:
		ch.handleBudget(update, state, userID)
	case 4:
		ch.handleTravelersCount(update, state, userID)
	case 5:
		ch.handleTravelDates(update, state, userID)
	case 6:
		ch.handleDuration(update, state, userID)
	case 7:
		ch.handleAccommodationType(update, state, userID)
	case 8:
		ch.handleSpecialRequirements(update, state, userID)
	case 9:
		ch.handleConfirmation(update, state, userID)
	default:
		logrus.WithField("user_id", userID).Warn("Неизвестный шаг, сброс состояния")
		ch.resetUserState(userID)
	}
}

func (ch *ConversationHandler) handleDestinationType(update tgbotapi.Update, state *models.TravelPreferences, userID int64) {
	state.DestinationType = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, 2)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`*Шаг 2 из 8:*
В какие страны или направления вы хотели бы поехать?
(Можно перечислить несколько через запятую)`)
	msg.ParseMode = "Markdown"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleCountries(update tgbotapi.Update, state *models.TravelPreferences, userID int64) {
	state.Countries = utils.ValidateCountries(update.Message.Text)
	ch.commandHandler.UpdateUserStep(userID, 3)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`*Шаг 3 из 8:*
Какой ориентировочный бюджет на *одного человека* (в рублях или валюте)?
(например: 50000 руб, 1500$)`)
	msg.ParseMode = "Markdown"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleBudget(update tgbotapi.Update, state *models.TravelPreferences, userID int64) {
	if !utils.ValidateBudget(update.Message.Text) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Пожалуйста, укажите бюджет с цифрами.\n"+
				"Например: *50000 руб* или *1500$*")
		msg.ParseMode = "Markdown"
		ch.commandHandler.bot.Send(msg)
		return
	}

	state.BudgetPerPerson = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, 4)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`*Шаг 4 из 8:*
Сколько человек планируют путешествие?
(Включая детей, укажите возраст детей, если есть)`)
	msg.ParseMode = "Markdown"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleTravelersCount(update tgbotapi.Update, state *models.TravelPreferences, userID int64) {
	state.TravelersCount = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, 5)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`*Шаг 5 из 8:*
На какие даты или период планируется поездка?
(например: *июль 2024*, *10-25 августа*, *на новый год*)`)
	msg.ParseMode = "Markdown"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleTravelDates(update tgbotapi.Update, state *models.TravelPreferences, userID int64) {
	state.TravelDates = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, 6)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`*Шаг 6 из 8:*
Какова желаемая продолжительность поездки?
(например: *7-10 дней*, *2 недели*, *выходные*)`)
	msg.ParseMode = "Markdown"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleDuration(update tgbotapi.Update, state *models.TravelPreferences, userID int64) {
	state.Duration = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, 7)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`*Шаг 7 из 8:*
Какой тип проживания предпочитаете?
(например: *отель 5*, *апартаменты*, *вилла*, *хостел*, *все включено*)`)
	msg.ParseMode = "Markdown"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleAccommodationType(update tgbotapi.Update, state *models.TravelPreferences, userID int64) {
	state.AccommodationType = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, 8)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`*Шаг 8 из 8:*
Есть ли особые пожелания или требования?
(например: *виза*, *питание*, *трансфер*, *доступная среда*)
Если нет, напишите 'нет' или 'нет особых'`)
	msg.ParseMode = "Markdown"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleSpecialRequirements(update tgbotapi.Update, state *models.TravelPreferences, userID int64) {
	state.SpecialRequirements = update.Message.Text
	state.CreatedAt = time.Now()
	ch.commandHandler.UpdateUserStep(userID, 9)

	// Формируем информацию о пользователе
	userInfo := models.UserInfo{
		ID:        update.Message.From.ID,
		FirstName: update.Message.From.FirstName,
		LastName:  update.Message.From.LastName,
		Username:  update.Message.From.UserName,
	}

	// Создаем превью
	preview := state.ToFormattedString(userInfo)

	// Сохраняем превью для подтверждения (в реальном приложении нужно хранить в состоянии)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf(`*Превью вашей заявки:*
%s

Всё верно? Отправьте *'да'* для подтверждения или *'нет'* для перезаполнения.`, preview))
	msg.ParseMode = "Markdown"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleConfirmation(update tgbotapi.Update, state *models.TravelPreferences, userID int64) {
	logrus.WithFields(logrus.Fields{
		"user_id":      userID,
		"answer":       update.Message.Text,
		"state_exists": state != nil,
	}).Info("Обработка подтверждения")

	answer := strings.ToLower(update.Message.Text)

	if strings.Contains(answer, "да") || strings.Contains(answer, "yes") || answer == "ок" {
		logrus.WithField("user_id", userID).Info("Пользователь подтвердил заявку")

		// Отправляем менеджеру
		userInfo := models.UserInfo{
			ID:        update.Message.From.ID,
			FirstName: update.Message.From.FirstName,
			LastName:  update.Message.From.LastName,
			Username:  update.Message.From.UserName,
		}

		logrus.WithFields(logrus.Fields{
			"user_id":      userID,
			"user_info":    fmt.Sprintf("%+v", userInfo),
			"travel_prefs": fmt.Sprintf("%+v", *state),
		}).Info("Данные для отправки менеджеру")

		if err := ch.formService.SendToManager(*state, userInfo); err != nil {
			logrus.WithError(err).Error("Ошибка при отправке заявки менеджеру")

			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"Произошла ошибка при отправке заявки. Пожалуйста, попробуйте позже.")
			ch.commandHandler.bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				`✅ *Спасибо! Ваша заявка отправлена менеджеру.*

Наш специалист свяжется с вами в течение 24 часов для уточнения деталей и подбора лучших предложений.

Для оформления новой заявки нажмите /newrequest`)
			msg.ParseMode = "Markdown"
			ch.commandHandler.bot.Send(msg)

			logrus.WithFields(logrus.Fields{
				"user_id":  userID,
				"username": userInfo.Username,
			}).Info("Заявка успешно отправлена")
		}

		ch.resetUserState(userID)

	} else if strings.Contains(answer, "нет") || strings.Contains(answer, "no") {
		logrus.WithField("user_id", userID).Info("Пользователь отказался от заявки")
		// Начинаем заново
		ch.resetUserState(userID)
		ch.commandHandler.HandleNewRequest(update)

	} else {
		logrus.WithField("user_id", userID).Warn("Непонятный ответ от пользователя")
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Пожалуйста, ответьте *'да'* для подтверждения или *'нет'* для перезаполнения.")
		msg.ParseMode = "Markdown"
		ch.commandHandler.bot.Send(msg)
	}
}

func (ch *ConversationHandler) resetUserState(userID int64) {
	delete(ch.commandHandler.userStates, userID)
	delete(ch.commandHandler.userStep, userID)
}
