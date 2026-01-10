package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	BotToken      string
	ManagerChatID int64
	DebugMode     bool
}

var AppConfig Config

func Load() error {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Файл .env не найден, используются переменные окружения")
	}

	AppConfig = Config{
		BotToken:      getEnv("BOT_TOKEN", ""),
		ManagerChatID: getEnvAsInt64("MANAGER_CHAT_ID", 0),
		DebugMode:     getEnvAsBool("DEBUG_MODE", true), // Включаем debug по умолчанию
	}

	logrus.Infof("Загружена конфигурация: ManagerChatID=%d, DebugMode=%v",
		AppConfig.ManagerChatID, AppConfig.DebugMode)

	return validate()
}

func validate() error {
	if AppConfig.BotToken == "" {
		logrus.Fatal("BOT_TOKEN не установлен")
	}

	if AppConfig.ManagerChatID == 0 {
		logrus.Error("MANAGER_CHAT_ID не установлен или равен 0. Заявки не будут пересылаться!")
	} else {
		logrus.Infof("MANAGER_CHAT_ID установлен: %d", AppConfig.ManagerChatID)
	}

	return nil
}

// Вспомогательные функции для чтения переменных окружения
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	logrus.Warnf("Переменная %s не найдена, используется значение по умолчанию: %s", key, defaultValue)
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
		logrus.Errorf("Не удалось преобразовать %s=%s в число", key, value)
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
