package main

import (
	"time"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/bot"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/dbconnection"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/joho/godotenv"
)

func main() {
	time.Local = time.UTC
	logger := log.GetLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	dbconnection.GetPoolConnections()

	logger.Info("All Rigth!")

	bot, err := bot.New()

	if err != nil {
		logger.Fatal(err)
	}

	bot.StartPolling()
}
