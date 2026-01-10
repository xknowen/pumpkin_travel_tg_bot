package handlers

import (
	"pumpkin_travel_tg_bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type CommandHandler struct {
	bot        *tgbotapi.BotAPI
	userStates map[int64]*models.TravelPreferences
	userStep   map[int64]int
}

func NewCommandHandler(bot *tgbotapi.BotAPI) *CommandHandler {
	return &CommandHandler{
		bot:        bot,
		userStates: make(map[int64]*models.TravelPreferences),
		userStep:   make(map[int64]int),
	}
}

func (ch *CommandHandler) HandleStart(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `🎉 *Добро пожаловать в TravelBot!*
Я помогу вам подобрать идеальное путешествие.

*Доступные команды:*
/newrequest - Начать оформление новой заявки
/help - Получить справку
/cancel - Отменить текущий диалог

Просто нажмите /newrequest, чтобы начать!`)
	msg.ParseMode = "Markdown"

	ch.bot.Send(msg)
	logrus.WithFields(logrus.Fields{
		"user_id":  update.Message.From.ID,
		"username": update.Message.From.UserName,
	}).Info("Пользователь запустил бота")
}

func (ch *CommandHandler) HandleHelp(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `*Помощь по боту*

Этот бот собирает ваши пожелания к путешествию и передает их нашему менеджеру.

*Как это работает:*
1. Нажмите /newrequest
2. Ответьте на вопросы о типе отдыха, бюджете, датах и т.д.
3. После заполнения всех данных заявка автоматически отправится менеджеру.
4. Менеджер свяжется с вами в течение 24 часов.

Вы можете прервать заполнение заявки командой /cancel.`)
	msg.ParseMode = "Markdown"

	ch.bot.Send(msg)
}

func (ch *CommandHandler) HandleCancel(update tgbotapi.Update) {
	userID := update.Message.From.ID

	if _, exists := ch.userStates[userID]; exists {
		delete(ch.userStates, userID)
		delete(ch.userStep, userID)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Диалог прерван. Ваши данные не сохранены.\n"+
			"Чтобы начать заново, нажмите /newrequest")

	ch.bot.Send(msg)
	logrus.WithField("user_id", userID).Info("Пользователь прервал диалог")
}

func (ch *CommandHandler) HandleNewRequest(update tgbotapi.Update) {
	userID := update.Message.From.ID

	// Инициализируем состояние пользователя
	ch.userStates[userID] = &models.TravelPreferences{}
	ch.userStep[userID] = 1

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`Отлично! Давайте подберем для вас идеальное путешествие. 🧳
Я задам несколько вопросов, это займет 2-3 минуты.

*Шаг 1 из 8:*
Какой тип отдыха вас интересует?
(например: *пляжный*, *экскурсионный*, *горнолыжный*, *гастрономический*)`)
	msg.ParseMode = "Markdown"

	ch.bot.Send(msg)
	logrus.WithField("user_id", userID).Info("Начался новый диалог с пользователем")
}

// Получение текущего состояния пользователя
func (ch *CommandHandler) GetUserState(userID int64) (*models.TravelPreferences, int, bool) {
	state, stateExists := ch.userStates[userID]
	step, stepExists := ch.userStep[userID]

	if !stateExists || !stepExists {
		return nil, 0, false
	}

	return state, step, true
}

// Обновление шага пользователя
func (ch *CommandHandler) UpdateUserStep(userID int64, step int) {
	ch.userStep[userID] = step
}
