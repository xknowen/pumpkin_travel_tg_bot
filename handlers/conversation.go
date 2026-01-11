package handlers

import (
	"fmt"
	"pumpkin_travel_tg_bot/models"
	"pumpkin_travel_tg_bot/services"
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

const (
	STEP_DESTINATION = iota + 1
	STEP_DEPARTURE_CITY
	STEP_TRAVEL_DATES
	STEP_DURATION
	STEP_TRAVELERS
	STEP_CHILD_AGE
	STEP_BUDGET
	STEP_VACATION_TYPE
	STEP_HOTEL_LEVEL
	STEP_MEAL_PLAN
	STEP_IMPORTANT_FACTORS
	STEP_CONFIRMATION
)

func (ch *ConversationHandler) HandleMessage(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		userID := update.CallbackQuery.From.ID
		state, step, exists := ch.commandHandler.GetUserState(userID)
		if exists && step == STEP_HOTEL_LEVEL {
			ch.handleHotelLevel(update, state, userID)
		}
		return
	}

	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID

	state, step, exists := ch.commandHandler.GetUserState(userID)
	if !exists {
		ch.commandHandler.HandleHelp(update)
		return
	}

	switch step {
	case STEP_DESTINATION:
		ch.handleDestination(update, state, userID)
	case STEP_DEPARTURE_CITY:
		ch.handleDepartureCity(update, state, userID)
	case STEP_TRAVEL_DATES:
		ch.handleTravelDates(update, state, userID)
	case STEP_DURATION:
		ch.handleDuration(update, state, userID)
	case STEP_TRAVELERS:
		ch.handleTravelers(update, state, userID)
	case STEP_CHILD_AGE:
		ch.handleChildAge(update, state, userID)
	case STEP_BUDGET:
		ch.handleBudget(update, state, userID)
	case STEP_VACATION_TYPE:
		ch.handleVacationType(update, state, userID)
	case STEP_HOTEL_LEVEL:
		ch.handleHotelLevel(update, state, userID)
	case STEP_MEAL_PLAN:
		ch.handleMealPlan(update, state, userID)
	case STEP_IMPORTANT_FACTORS:
		ch.handleImportantFactors(update, state, userID)
	case STEP_CONFIRMATION:
		ch.handleConfirmation(update, state, userID)
	default:
		ch.resetUserState(userID)
	}
}

func (ch *ConversationHandler) handleDestination(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	state.Destination = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, STEP_DEPARTURE_CITY)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`2️⃣
<b>Из какого города планируется вылет?</b>
(Напишите ваш город или из которого хотите вылететь)

<code>Например: Москва, Краснодар или Сочи</code>`)
	msg.ParseMode = "HTML"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleDepartureCity(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	state.DepartureCity = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, STEP_TRAVEL_DATES)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`3️⃣
<b>Желаемые даты поездки</b>
(Напишите точные даты или примерные)

<code>Например:
10–20 мая
Июнь
Любые даты февраля
Самые бюджетные на следующий месяц</code>`)
	msg.ParseMode = "HTML"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleTravelDates(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	state.TravelDates = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, STEP_DURATION)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`4️⃣
<b>Сколько дней планируете отдых?</b>
(Напишите точное или примерное количество)

<code>Например: 3 дня / неделя / 10–14 дней</code>`)
	msg.ParseMode = "HTML"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleDuration(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	state.Duration = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, STEP_TRAVELERS)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`5️⃣
<b>Сколько человек летит?</b>
(Напишите количество туристов)

<code>Например:
2 взрослых
2 взрослых + 1 ребёнок
1 взрослый</code>`)
	msg.ParseMode = "HTML"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleTravelers(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	state.Travelers = update.Message.Text

	answer := strings.ToLower(update.Message.Text)
	if strings.Contains(answer, "ребен") || strings.Contains(answer, "дет") {
		ch.commandHandler.UpdateUserStep(userID, STEP_CHILD_AGE)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			`<b>Сколько лет ребенку?</b>
(Напишите возраст)

<code>Например: 3 года / 5 / 12 лет</code>`)
		msg.ParseMode = "HTML"
		ch.commandHandler.bot.Send(msg)
	} else {
		state.ChildAge = "Нет детей"
		ch.commandHandler.UpdateUserStep(userID, STEP_BUDGET)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			`6️⃣
<b>Бюджет на всех (перелёт + проживание)</b>
(Напишите планируемый бюджет)

<code>Например:
до 80 000 ₽
200–250 тыс.
Без строгих рамок</code>`)
		msg.ParseMode = "HTML"
		ch.commandHandler.bot.Send(msg)
	}
}

func (ch *ConversationHandler) handleChildAge(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	state.ChildAge = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, STEP_BUDGET)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`6️⃣
<b>Бюджет на всех (перелёт + проживание)</b>
(Напишите планируемый бюджет)

<code>Например:
до 80 000 ₽
200–250 тыс.
Без строгих рамок</code>`)
	msg.ParseMode = "HTML"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleBudget(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	state.Budget = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, STEP_VACATION_TYPE)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`7️⃣
<b>Какой отдых вы хотите?</b>
(Напишите все пожелания по отдыху)

<code>Например:
Пляжный
Пляж + экскурсии + все включено
Активный без детей
Спокойный / релакс
С детьми</code>`)
	msg.ParseMode = "HTML"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleVacationType(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	state.VacationType = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, STEP_HOTEL_LEVEL)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`8️⃣
<b>Какой уровень отеля рассматриваете?</b>

Выберите вариант ниже или напишите свой:`)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3★", "hotel_3"),
			tgbotapi.NewInlineKeyboardButtonData("4★", "hotel_4"),
			tgbotapi.NewInlineKeyboardButtonData("5★", "hotel_5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Любой уровень", "hotel_any"),
			tgbotapi.NewInlineKeyboardButtonData("Не имеет значения", "hotel_no_matter"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3★ или 4★", "hotel_3_4"),
			tgbotapi.NewInlineKeyboardButtonData("4★ или 5★", "hotel_4_5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отель 16+", "hotel_16"),
			tgbotapi.NewInlineKeyboardButtonData("Отель 18+", "hotel_18"),
		),
	)

	msg.ReplyMarkup = keyboard
	msg.ParseMode = "HTML"

	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleHotelLevel(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		if _, err := ch.commandHandler.bot.Request(callback); err != nil {
			logrus.Error("Ошибка отправки callback:", err)
		}

		callbackData := update.CallbackQuery.Data
		var hotelLevelText string

		switch callbackData {
		case "hotel_3":
			hotelLevelText = "3★"
		case "hotel_4":
			hotelLevelText = "4★"
		case "hotel_5":
			hotelLevelText = "5★"
		case "hotel_any":
			hotelLevelText = "Любой уровень"
		case "hotel_no_matter":
			hotelLevelText = "Не имеет значения"
		case "hotel_3_4":
			hotelLevelText = "3★ или 4★"
		case "hotel_4_5":
			hotelLevelText = "4★ или 5★"
		case "hotel_16":
			hotelLevelText = "Отель 16+"
		case "hotel_18":
			hotelLevelText = "Отель 18+"
		default:
			hotelLevelText = "Не указано"
		}

		state.HotelLevel = hotelLevelText

		editMsg := tgbotapi.NewEditMessageText(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			fmt.Sprintf(`✅ <b>Выбрано:</b> %s

9️⃣
<b>Желаемый тип питания</b>

<code>Наример:
Завтрак
Обед
Завтрак + ужин
Всё включено
Без разницы</code>`, hotelLevelText),
		)
		editMsg.ParseMode = "HTML"
		editMsg.ReplyMarkup = nil

		ch.commandHandler.bot.Send(editMsg)

		ch.commandHandler.UpdateUserStep(userID, STEP_MEAL_PLAN)

	} else if update.Message != nil {
		state.HotelLevel = update.Message.Text
		ch.commandHandler.UpdateUserStep(userID, STEP_MEAL_PLAN)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			`9️⃣
<b>Желаемый тип питания</b>

<code>Наример:
Завтрак
Обед
Завтрак + ужин
Всё включено
Без разницы</code>`)
		msg.ParseMode = "HTML"
		ch.commandHandler.bot.Send(msg)
	}
}

func (ch *ConversationHandler) handleMealPlan(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	state.MealPlan = update.Message.Text
	ch.commandHandler.UpdateUserStep(userID, STEP_IMPORTANT_FACTORS)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`🔟
<b>Что для вас принципиально важно?</b>

<code>Например:
Первая линия
Песчаный пляж
Хороший Wi-Fi
Без пересадок
Свой бассейн</code>

<em>Если ничего не принципиально — напишите "нет"</em>`)
	msg.ParseMode = "HTML"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleImportantFactors(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	state.ImportantFactors = update.Message.Text
	state.CreatedAt = time.Now()
	ch.commandHandler.UpdateUserStep(userID, STEP_CONFIRMATION)

	preview := state.ToClientPreview()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf(`<b>✅ Все готово! Проверьте вашу заявку:</b>

%s

<b>Всё верно?</b> Отправьте <b>"да"</b> для подтверждения или <b>"нет"</b> для перезаполнения.`, preview))
	msg.ParseMode = "HTML"
	ch.commandHandler.bot.Send(msg)
}

func (ch *ConversationHandler) handleConfirmation(update tgbotapi.Update, state *models.TravelRequest, userID int64) {
	answer := strings.ToLower(update.Message.Text)

	if strings.Contains(answer, "да") || strings.Contains(answer, "yes") || answer == "ок" || answer == "подтверждаю" {
		userInfo := models.UserInfo{
			ID:        update.Message.From.ID,
			FirstName: update.Message.From.FirstName,
			LastName:  update.Message.From.LastName,
			Username:  update.Message.From.UserName,
		}

		if err := ch.formService.SendToManager(*state, userInfo); err != nil {
			logrus.WithError(err).Error("Ошибка при отправке заявки менеджеру")

			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"❌ Произошла ошибка при отправке заявки. Пожалуйста, попробуйте позже.")
			ch.commandHandler.bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				`✅ <b>Спасибо! Ваша заявка отправлена Ангелине.</b>

Ангелина свяжется с вами в ближайшее время для подбора лучших вариантов.

Для оформления новой заявки нажмите /newrequest`)
			msg.ParseMode = "HTML"
			ch.commandHandler.bot.Send(msg)

			logrus.WithFields(logrus.Fields{
				"user_id":  userID,
				"username": userInfo.Username,
			}).Info("Заявка успешно отправлена")
		}

		ch.resetUserState(userID)

	} else if strings.Contains(answer, "нет") || strings.Contains(answer, "no") {
		ch.resetUserState(userID)
		ch.commandHandler.HandleNewRequest(update)

	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Пожалуйста, ответьте <b>\"да\"</b> для подтверждения или <b>\"нет\"</b> для перезаполнения.")
		msg.ParseMode = "HTML"
		ch.commandHandler.bot.Send(msg)
	}
}

func (ch *ConversationHandler) resetUserState(userID int64) {
	delete(ch.commandHandler.userStates, userID)
	delete(ch.commandHandler.userStep, userID)
}
