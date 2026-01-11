package bot

import (
	"fmt"
	"pumpkin_travel_tg_bot/config"
	"pumpkin_travel_tg_bot/handlers"
	"pumpkin_travel_tg_bot/models"
	"pumpkin_travel_tg_bot/services"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type TravelBot struct {
	botAPI         *tgbotapi.BotAPI
	commandHandler *handlers.CommandHandler
	convHandler    *handlers.ConversationHandler
	formService    *services.FormService
}

func NewTravelBot() (*TravelBot, error) {
	if err := config.Load(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	botAPI, err := tgbotapi.NewBotAPI(config.AppConfig.BotToken)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания бота: %w", err)
	}

	botAPI.Debug = config.AppConfig.DebugMode

	logrus.Infof("Авторизован как %s", botAPI.Self.UserName)
	logrus.Infof("ID бота: %d", botAPI.Self.ID)

	formService := services.NewFormService(botAPI)
	commandHandler := handlers.NewCommandHandler(botAPI)
	convHandler := handlers.NewConversationHandler(commandHandler, formService)

	return &TravelBot{
		botAPI:         botAPI,
		commandHandler: commandHandler,
		convHandler:    convHandler,
		formService:    formService,
	}, nil
}

func (tb *TravelBot) Start() error {
	logrus.Info("Бот запускается...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tb.botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			tb.handleCallbackQuery(update)
			continue
		}

		if update.Message == nil {
			continue
		}

		logrus.WithFields(logrus.Fields{
			"user_id":  update.Message.From.ID,
			"username": update.Message.From.UserName,
			"text":     update.Message.Text,
			"chat_id":  update.Message.Chat.ID,
		}).Debug("Получено сообщение")

		if update.Message.IsCommand() {
			tb.handleCommand(update)
		} else {
			tb.convHandler.HandleMessage(update)
		}
	}

	return nil
}

func (tb *TravelBot) handleCallbackQuery(update tgbotapi.Update) {
	userID := update.CallbackQuery.From.ID

	logrus.WithFields(logrus.Fields{
		"user_id":       userID,
		"callback_data": update.CallbackQuery.Data,
	}).Info("Обработка callback query")

	_, step, exists := tb.commandHandler.GetUserState(userID)
	if !exists {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Диалог не активен. Начните заново /newrequest")
		tb.botAPI.Send(callback)
		return
	}

	if step == handlers.STEP_HOTEL_LEVEL {
		tb.convHandler.HandleMessage(update)
	} else {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Неверный шаг диалога")
		tb.botAPI.Send(callback)
	}
}

func (tb *TravelBot) handleCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		tb.commandHandler.HandleStart(update)
	case "help":
		tb.commandHandler.HandleHelp(update)
	case "newrequest":
		tb.commandHandler.HandleNewRequest(update)
	case "cancel":
		tb.commandHandler.HandleCancel(update)
	case "test":
		tb.handleTestCommand(update)
	case "config":
		tb.handleConfig(update)
	case "myid":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Ваш Chat ID: `%d`", update.Message.Chat.ID))
		msg.ParseMode = "Markdown"
		tb.botAPI.Send(msg)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Неизвестная команда. Используйте /help для списка команд")
		tb.botAPI.Send(msg)
	}
}

func (tb *TravelBot) handleTestCommand(update tgbotapi.Update) {
	logrus.Info("Вызвана тестовая команда /test")

	testRequest := models.TravelRequest{
		Destination:      "Тестовая страна",
		DepartureCity:    "Тестовый город",
		TravelDates:      "тест",
		Duration:         "7 дней",
		Travelers:        "2 взрослых",
		ChildAge:         "Нет детей",
		Budget:           "100000 ₽",
		VacationType:     "Пляжный",
		HotelLevel:       "4★",
		MealPlan:         "All Inclusive",
		ImportantFactors: "тест",
		CreatedAt:        time.Now(),
	}

	userInfo := models.UserInfo{
		ID:        update.Message.From.ID,
		FirstName: update.Message.From.FirstName,
		LastName:  update.Message.From.LastName,
		Username:  update.Message.From.UserName,
	}

	err := tb.formService.SendToManager(testRequest, userInfo)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if err != nil {
		msg.Text = fmt.Sprintf("❌ Ошибка отправки тестовой заявки: %v", err)
	} else {
		msg.Text = "✅ Тестовая заявка отправлена менеджеру. Проверьте, получил ли он её."
	}

	tb.botAPI.Send(msg)
}

func (tb *TravelBot) handleConfig(update tgbotapi.Update) {
	configInfo := fmt.Sprintf(
		"📋 *Конфигурация бота:*\n"+
			"• Имя бота: %s\n"+
			"• ID бота: %d\n"+
			"• ManagerChatID: `%d`\n"+
			"• Debug mode: %v\n\n"+
			"Для теста отправки используйте /test",
		tb.botAPI.Self.UserName,
		tb.botAPI.Self.ID,
		config.AppConfig.ManagerChatID,
		config.AppConfig.DebugMode)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, configInfo)
	msg.ParseMode = "Markdown"
	tb.botAPI.Send(msg)
}
