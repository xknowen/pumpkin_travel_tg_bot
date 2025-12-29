package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type config struct {
	BotToken      string
	ManagerChatID int64
	DebugMode     Bool
}

var AppConfig Config

func Load() error {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Файл .env не найден, используются переменные окружения")
	}

	AppConfig = Config{
		BotToken:      getEnv("BOT_TOKEN", ""),
		ManagerChatID: getEnv("MANAGER_CHAT_ID", 0),
		DebugMode:     getEnv("DEBUG_MODE", false),
	}
	return validate
}

func validate() error {
	if AppConfig.BotToken == "" {
		logrus.Fatal("BOT_TOKEN не установлен")
	}
	if AppConfig.ManagerChatID == 0 {
		logrus.Warn("MANAGER_CHAT_ID не установлен. Заявки не будут приниматься")
	}

	return nil
}
