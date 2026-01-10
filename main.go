package main

import (
	"pumpkin_travel_tg_bot/internal/bot"

	"github.com/sirupsen/logrus"
)

func main() {
	// Настраиваем логирование
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	// Создаем и запускаем бота
	travelBot, err := bot.NewTravelBot()
	if err != nil {
		logrus.Fatalf("Ошибка создания бота: %v", err)
	}

	// Запускаем бота
	if err := travelBot.Start(); err != nil {
		logrus.Fatalf("Ошибка запуска бота: %v", err)
	}
}
