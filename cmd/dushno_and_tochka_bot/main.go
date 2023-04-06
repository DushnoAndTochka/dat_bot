package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/bot"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/dbconnection"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/storages"
	"github.com/joho/godotenv"
)

// Точка запуска бота. Инициализирует все основные куски проекта и вызвывает бот поллинг.
func main() {
	time.Local = time.UTC
	logger := log.GetLogger()

	err := godotenv.Load()

	if err != nil {
		logger.Error(err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	storages.NewStorage(ctx)

	pool := dbconnection.GetPoolConnections()
	if pool == nil {
		logger.Fatal("Pool is not init")
	}

	logger.Info("All Rigth!")

	tgBot, err := bot.New()

	if err != nil {
		logger.Fatal(err)
	}

	bot.StartPolling(tgBot, ctx)
}
