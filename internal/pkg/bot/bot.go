package bot

import (
	"context"
	"os"
	"sync"

	"github.com/mymmrac/telego"

	th "github.com/mymmrac/telego/telegohandler"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/custompredicates"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/dbconnection"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/eventprocessor"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
)

// Инициализатор Бота. Ожидает что в переменных окржения имеет переменная TELEGRAM_BOT_TOKEN,
// в которой указан токен. Если данной переменной окружения не будет, то вызовется паника
func New() (*telego.Bot, error) {
	logger := log.GetLogger()
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if telegramBotToken == "" {
		logger.Fatal("Telegram Bot Token not found. Please specify TELEGRAM_BOT_TOKEN env.")
	}

	newApiBot, err := telego.NewBot(telegramBotToken, telego.WithLogger(logger))

	if err != nil {
		return nil, err
	}

	return newApiBot, nil

}

func StartPolling(bot *telego.Bot, ctx context.Context) {
	logger := log.GetLogger()
	var wg sync.WaitGroup

	updates, _ := bot.UpdatesViaLongPolling(nil)
	dbpool := dbconnection.GetPoolConnections()

	bh, _ := th.NewBotHandler(bot, updates)
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		<-ctx.Done()

		logger.Info("Polling is stoping")
		bot.StopLongPolling()
		bh.Stop()

		logger.Info("Long polling stoped")

		wg.Done()
	}(&wg)

	defer bh.Stop()
	defer bot.StopLongPolling()
	defer dbpool.Close()

	bh.Use(
		func(next th.Handler) th.Handler {
			return func(bot *telego.Bot, update telego.Update) {
				// midleware для оборачивания обработки в горутину.
				go func() {
					defer func() {
						if r := recover(); r != nil {
							logger.Error("panic recovered: %s", r)
						}
					}()
					next(bot, update)
				}()
			}
		},
	)

	initialHandlers(bh)
	bh.Start()
	wg.Wait()
	logger.Info("Long polling is done")
}

// инициализируем все обработчики
func initialHandlers(bh *th.BotHandler) {
	bh.Handle(eventprocessor.ProcessStartComand, th.CommandEqual("start"))
	bh.Handle(eventprocessor.ProcessHelpComand, th.CommandEqual("help"))
	bh.Handle(eventprocessor.ProcessProposeProblemFromMessage, th.CommandEqual("suggest_problem"))
	bh.Handle(eventprocessor.ProcessShowAllProposeProblems, th.CommandEqual("show_top_suggestions"))
	bh.Handle(eventprocessor.ProcessShowMyProposeProblem, th.CommandEqual("show_my_suggestion"))
	bh.Handle(eventprocessor.ProcessProposeProblemFromMessage, th.TextEqual("Хочу предложить задачу"))
	bh.Handle(eventprocessor.ProcessGetLinkFromReply, custompredicates.IsNewProposeTask)
	bh.Handle(eventprocessor.ProcessNotSupportedComandsComand, th.AnyCommand())
}
