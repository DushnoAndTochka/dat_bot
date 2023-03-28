package bot

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymmrac/telego"

	th "github.com/mymmrac/telego/telegohandler"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/custompredicates"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/dbconnection"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/eventprocessor"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
)

func New() (*Bot, error) {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if telegramBotToken == "" {
		err := errors.New("Telegram Bot Token not found. Please specify TELEGRAM_BOT_TOKEN env.")
		return nil, err
	}

	logger := log.GetLogger()

	newApiBot, err := telego.NewBot(telegramBotToken, telego.WithLogger(logger))

	if err != nil {
		return nil, err
	}

	dbConnection := dbconnection.GetPoolConnections()
	botHandler := &Bot{
		bot:          newApiBot,
		dbConnection: dbConnection,
	}

	return botHandler, nil

}

type Bot struct {
	bot          *telego.Bot
	dbConnection *pgxpool.Pool
}

func (b *Bot) StartPolling() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	logger := log.GetLogger()
	done := make(chan struct{}, 1)

	updates, _ := b.bot.UpdatesViaLongPolling(nil)

	bh, _ := th.NewBotHandler(b.bot, updates)

	go func() {
		<-sigs

		logger.Info("Polling is stoping")
		b.bot.StopLongPolling()
		bh.Stop()
		logger.Info("Long polling stoped")

		done <- struct{}{}
	}()

	defer bh.Stop()
	defer b.bot.StopLongPolling()

	bh.Use(
		func(next th.Handler) th.Handler {
			return func(bot *telego.Bot, update telego.Update) {
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

	bh.Handle(eventprocessor.ProcessStartComand, th.CommandEqual("start"))
	bh.Handle(eventprocessor.ProcessHelpComand, th.CommandEqual("help"))
	bh.Handle(eventprocessor.ProcessProposeProblemFromMessage, th.CommandEqual("propose_problem"))
	bh.Handle(eventprocessor.ProcessShowAllProposeProblems, th.CommandEqual("show_top_propose"))
	bh.Handle(eventprocessor.ProcessShowMyProposeProblem, th.CommandEqual("show_my_propose"))
	bh.Handle(eventprocessor.ProcessProposeProblemFromMessage, th.TextEqual("Хочу предложить задачу"))
	bh.Handle(eventprocessor.ProcessGetLinkFromReply, custompredicates.IsNewProposeTask)
	bh.Handle(eventprocessor.ProcessNotSupportedComandsComand, th.AnyCommand())

	bh.Start()

	<-done
	logger.Info("Long polling is done")
}
