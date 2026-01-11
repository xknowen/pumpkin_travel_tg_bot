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

func (fs *FormService) SendToManager(request models.TravelRequest, userInfo models.UserInfo) error {
	if config.AppConfig.ManagerChatID == 0 {
		return fmt.Errorf("MANAGER_CHAT_ID не задан")
	}

	messageText := request.ToFormattedString(userInfo)

	msg := tgbotapi.NewMessage(config.AppConfig.ManagerChatID, messageText)
	msg.ParseMode = "HTML"

	if _, err := fs.bot.Send(msg); err != nil {
		logrus.WithError(err).Error("Ошибка при отправке заявки менеджеру")
		return err
	}

	logrus.Info("✅ Заявка успешно отправлена менеджеру")
	return nil
}
