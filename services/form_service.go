package services

import (
	"fmt"
	"pumpkin_travel_tg_bot/config"
	"pumpkin_travel_tg_bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type FormService struct {
	bot *tgbotapi.BotAPI
}

func NewFormService(bot *tgbotapi.BotAPI) *FormService {
	return &FormService{bot: bot}
}

func (fs *FormService) SendToManager(preferences models.TravelPreferences, userInfo models.UserInfo) error {
	logrus.Infof("Попытка отправить заявку менеджеру. ManagerChatID: %d", config.AppConfig.ManagerChatID)

	if config.AppConfig.ManagerChatID == 0 {
		logrus.Error("MANAGER_CHAT_ID не задан или равен 0")
		return fmt.Errorf("MANAGER_CHAT_ID не задан")
	}

	messageText := preferences.ToFormattedString(userInfo)

	logrus.Infof("Текст заявки для менеджера:\n%s", messageText)

	msg := tgbotapi.NewMessage(config.AppConfig.ManagerChatID, messageText)
	msg.ParseMode = "Markdown"

	logrus.Info("Отправка сообщения менеджеру...")

	if _, err := fs.bot.Send(msg); err != nil {
		logrus.WithError(err).Error("Ошибка при отправке заявки менеджеру")

		// Проверим конкретную ошибку
		if err.Error() == "Forbidden: bot was blocked by the user" {
			logrus.Error("Бот заблокирован пользователем (менеджером)")
		} else if err.Error() == "Bad Request: chat not found" {
			logrus.Error("Чат с менеджером не найден. Проверьте ManagerChatID")
		}

		return err
	}

	logrus.Info("✅ Заявка успешно отправлена менеджеру")
	return nil
}
