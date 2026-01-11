package main

import (
	"pumpkin_travel_tg_bot/internal/bot"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	travelBot, err := bot.NewTravelBot()
	if err != nil {
		logrus.Fatalf("Ошибка создания бота: %v", err)
	}

	if err := travelBot.Start(); err != nil {
		logrus.Fatalf("Ошибка запуска бота: %v", err)
	}
}
