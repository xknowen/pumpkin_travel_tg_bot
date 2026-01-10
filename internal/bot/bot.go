package bot

import (
	"fmt"
	"pumpkin_travel_tg_bot/config"
	"pumpkin_travel_tg_bot/handlers"
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
	// Загружаем конфигурацию
	if err := config.Load(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	// Создаем экземпляр бота
	botAPI, err := tgbotapi.NewBotAPI(config.AppConfig.BotToken)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания бота: %w", err)
	}

	botAPI.Debug = config.AppConfig.DebugMode

	logrus.Infof("Авторизован как %s", botAPI.Self.UserName)
	logrus.Infof("ID бота: %d", botAPI.Self.ID)

	// Создаем сервисы и обработчики
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

	// Настраиваем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tb.botAPI.GetUpdatesChan(u)

	// Обрабатываем обновления
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Логируем входящее сообщение
		logrus.WithFields(logrus.Fields{
			"user_id":  update.Message.From.ID,
			"username": update.Message.From.UserName,
			"text":     update.Message.Text,
			"chat_id":  update.Message.Chat.ID,
		}).Debug("Получено сообщение")

		// Обрабатываем команды
		if update.Message.IsCommand() {
			tb.handleCommand(update)
		} else {
			// Обрабатываем обычные сообщения (диалог)
			tb.convHandler.HandleMessage(update)
		}
	}

	return nil
}

func (tb *TravelBot) handleCommand(update tgbotapi.Update) {
	logrus.Infof("Обработка команды: %s от пользователя %d",
		update.Message.Command(), update.Message.From.ID)

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
		tb.handleConfig(update) // Новая команда
	case "myid":
		// Добавьте эту команду если нужно
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
	logrus.Info("=== ВЫЗВАНА ТЕСТОВАЯ КОМАНДА /test ===")

	// Простое тестовое сообщение
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🔄 Отправляю тестовое сообщение менеджеру...")
	tb.botAPI.Send(msg)

	// Пробуем отправить простое сообщение менеджеру
	testMsg := tgbotapi.NewMessage(config.AppConfig.ManagerChatID,
		"🔄 *Тестовая заявка от бота*\n"+
			"Если вы видите это сообщение, значит бот может отправлять заявки менеджеру!\n"+
			"Время: "+time.Now().Format("15:04 02.01.2006"))
	testMsg.ParseMode = "Markdown"

	logrus.Infof("Отправляю тестовое сообщение в чат: %d", config.AppConfig.ManagerChatID)

	if _, err := tb.botAPI.Send(testMsg); err != nil {
		logrus.Errorf("❌ Ошибка отправки тестового сообщения: %v", err)

		// Отправляем пользователю информацию об ошибке
		errorMsg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка отправки: %v\n\nПроверьте:\n1. Правильность ManagerChatID в .env\n2. Что менеджер не заблокировал бота\n3. Что менеджер писал боту хоть раз", err))
		tb.botAPI.Send(errorMsg)
	} else {
		logrus.Info("✅ Тестовое сообщение отправлено успешно!")

		successMsg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"✅ Тестовое сообщение отправлено менеджеру!\n"+
				"Проверьте, получил ли менеджер (пользователь с ID 1990105992) это сообщение.\n"+
				"Если не получил, возможно, он заблокировал бота.")
		tb.botAPI.Send(successMsg)
	}
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
