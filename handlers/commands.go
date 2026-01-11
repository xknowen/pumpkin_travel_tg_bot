package handlers

import (
	"pumpkin_travel_tg_bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type CommandHandler struct {
	bot        *tgbotapi.BotAPI
	userStates map[int64]*models.TravelRequest
	userStep   map[int64]int
}

func NewCommandHandler(bot *tgbotapi.BotAPI) *CommandHandler {
	return &CommandHandler{
		bot:        bot,
		userStates: make(map[int64]*models.TravelRequest),
		userStep:   make(map[int64]int),
	}
}

func (ch *CommandHandler) HandleStart(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`🤍 <b>Привет!</b>
Я — помогаю подобрать путешествия без хлопот и лишней суеты ✈️

Подбираю туры под конкретные даты, бюджет и формат отдыха — так, как подбирала бы для себя.

Ответьте на 10 коротких вопросов, и я предложу подходящие варианты 🌴

<b>Доступные команды:</b>
/newrequest — Начать оформление новой заявки
/help — Получить справку
/cancel — Отменить текущий диалог

Просто нажмите /newrequest, чтобы начать!`)
	msg.ParseMode = "HTML"

	ch.bot.Send(msg)
	logrus.WithFields(logrus.Fields{
		"user_id":  update.Message.From.ID,
		"username": update.Message.From.UserName,
	}).Info("Пользователь запустил бота")
}

func (ch *CommandHandler) HandleHelp(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`<b>Помощь по боту</b>

Этот бот собирает ваши пожелания к путешествию и передает их Ангелине — специалисту по подбору туров.

<b>Как это работает:</b>
1. Нажмите /newrequest
2. Ответьте на 10 вопросов о вашем путешествии
3. После заполнения всех данных заявка автоматически отправится
4. Ангелина свяжется с вами в ближайшее время с подбором вариантов

Вы можете прервать заполнение заявки командой /cancel в любой момент.`)
	msg.ParseMode = "HTML"

	ch.bot.Send(msg)
}

func (ch *CommandHandler) HandleCancel(update tgbotapi.Update) {
	userID := update.Message.From.ID

	if _, exists := ch.userStates[userID]; exists {
		delete(ch.userStates, userID)
		delete(ch.userStep, userID)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"❌ Диалог прерван. Ваши данные не сохранены.\n\nЧтобы начать заново, нажмите /newrequest")
	msg.ParseMode = "HTML"

	ch.bot.Send(msg)
	logrus.WithField("user_id", userID).Info("Пользователь прервал диалог")
}

func (ch *CommandHandler) HandleNewRequest(update tgbotapi.Update) {
	userID := update.Message.From.ID

	ch.userStates[userID] = &models.TravelRequest{}
	ch.userStep[userID] = STEP_DESTINATION

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		`🌴 <b>Отлично! Давайте подберем для вас идеальное путешествие.</b>

Я задам 10 вопросов, это займет 2-3 минуты.

1️⃣
<b>Куда планируете поездку?</b>
(Написать интересные вам направления)

<code>Пример: Турция / Россия / Пока не определились</code>

<em>Если нет конкретной страны — подберу варианты</em>`)
	msg.ParseMode = "HTML"

	ch.bot.Send(msg)
	logrus.WithField("user_id", userID).Info("Начался новый диалог с пользователем")
}

func (ch *CommandHandler) GetUserState(userID int64) (*models.TravelRequest, int, bool) {
	state, stateExists := ch.userStates[userID]
	step, stepExists := ch.userStep[userID]

	if !stateExists || !stepExists {
		return nil, 0, false
	}

	return state, step, true
}

func (ch *CommandHandler) UpdateUserStep(userID int64, step int) {
	ch.userStep[userID] = step
}
