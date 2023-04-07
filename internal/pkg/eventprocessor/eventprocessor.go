package eventprocessor

import (
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/customerrors"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/custompredicates"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/storages"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// Обработчик команды "/start". Пишет приветственное сообщение и добавляет новых пользователей в БД.
func ProcessStartComand(bot *telego.Bot, update telego.Update) {

	storage := storages.GetStorage()
	user := models.GetFromTg(&update)

	var message *telego.SendMessageParams
	if err := storage.UserGetOrCreate(user); err != nil {
		message = tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Привет. Чем могут быть полезен ?",
		)
	} else {
		message = tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("Привет, %s. Чем могут быть полезен ?", user.Name),
		)
	}

	_, _ = bot.SendMessage(message)
}

// Обработчик команды "/help". Пишет информационное сообщение о возможностях бота.
func ProcessHelpComand(bot *telego.Bot, update telego.Update) {

	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		`На данный момент я умею не так много.
		 - Через меня можно предложить задачу на разбор. /suggest_problem
		 - Посмотреть топ 10 желаемых задач от подписчиков. /show_top_suggestions
		 - Посмотреть какие задачи я уже предлагал. /show_my_suggestion
		 
		 На этом пока что все.`,
	)

	_, _ = bot.SendMessage(message)
}

// Обработчик не известных команд.
func ProcessNotSupportedComandsComand(bot *telego.Bot, update telego.Update) {
	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Неизвестная команда, используй /help.",
	).WithReplyMarkup(tu.ReplyKeyboardRemove())

	_, _ = bot.SendMessage(message)
}

// Обработчик предложенной задачи. Вычитывает ссылку из сообщения, валидирует ее и добавляет в БД.
func ProcessGetLinkFromReply(bot *telego.Bot, update telego.Update) {
	logger := log.GetLogger()
	answer := update.Message.Text
	storage := storages.GetStorage()

	user := models.GetFromTg(&update)

	if err := storage.UserGetOrCreate(user); err != nil {
		logger.Debug("ProcessGetLinkFromReply: UserGetOrCreate failed: %w", err)
		sendErrorMessage(bot, &update, err)
		return
	}

	problem, err := models.NewProblemFromUrl(answer)

	var message *telego.SendMessageParams

	if err != nil {
		sendErrorMessage(bot, &update, err)
		return
	}

	if err := storage.ProblemGetOrCreate(problem); err != nil {
		logger.Debug("ProcessGetLinkFromReply: ProblemGetOrCreate failed: %w", err)
		sendErrorMessage(bot, &update, err)
		return
	}

	suggection := models.NewSuggestion(user, problem)

	if err = storage.SuggestionCheckOrCreate(suggection); err != nil {
		logger.Debug("ProcessGetLinkFromReply: SuggestionCheckOrCreate failed: %w", err)
		sendErrorMessage(bot, &update, err)
		return
	}

	message = tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Спасибо",
	)
	_, _ = bot.SendMessage(message)
}

// Обработчик команды "/show_top_suggestions". Выбирает ТОП предложений из бд и формирует сообщение.
func ProcessShowAllProposeProblems(bot *telego.Bot, update telego.Update) {

	var message string
	var botMessage *telego.SendMessageParams

	logger := log.GetLogger()
	storage := storages.GetStorage()

	suggstions, err := storage.GetTopSuggestions()

	if err != nil {
		logger.Error("ProcessShowAllProposeProblems: %w", err)
		sendErrorMessage(bot, &update, err)
		return
	}

	if len(suggstions) == 0 {
		message =
			`
			На данный момент нет предложенных задач. Вы можете быть первым.
			/suggest_problem
			`
		botMessage = tu.Message(tu.ID(update.Message.Chat.ID), message)
	} else {
		// Формируем сообщение с ссылками на Problems.

		// инитим список сообщений
		var entityMessages []tu.MessageEntityCollection
		entityMessages = append(entityMessages, tu.Entity("ТОП предложениями являются:\n"))

		var problemUrl string

		for _, problem := range suggstions {
			problemUrl = problem.GetOriginalUrl()
			// собираем каждую строчку отдельно
			entityMessages = append(entityMessages, tu.Entity("\n- Задача "))
			if problemUrl != "" {
				entityMessages = append(entityMessages, tu.Entity(string(problem.Name)).TextLink(problemUrl))
			} else {
				entityMessages = append(entityMessages, tu.Entity(string(problem.Name)))
			}
			entityMessages = append(entityMessages, tu.Entity(fmt.Sprintf(" была предложена %v раз.\n", int(problem.CountSuggestions))))
		}
		// скармливаем собранный список сообщений и формируем целостное сообщение.
		botMessage = tu.MessageWithEntities(tu.ID(update.Message.Chat.ID), entityMessages...)

	}

	_, _ = bot.SendMessage(botMessage)
}

// Обработчик команды "/show_my_suggestion". Показывает предложения пользователя.
func ProcessShowMyProposeProblem(bot *telego.Bot, update telego.Update) {
	var botMessage *telego.SendMessageParams

	logger := log.GetLogger()
	user := models.GetFromTg(&update)
	storage := storages.GetStorage()

	if err := storage.UserGetOrCreate(user); err != nil {
		sendErrorMessage(bot, &update, err)
		return
	}

	userSuggestions, err := storage.GetUserSuggestion(user)

	if err != nil {
		logger.Error("ProcessShowMyProposeProblem: %w", err)
		sendErrorMessage(bot, &update, err)
		return
	}

	if userSuggestions != nil {
		// формируем итоговое сообщение с ссылками.
		var entityMessages []tu.MessageEntityCollection
		entityMessages = append(entityMessages, tu.Entity("Задачи, которые вы предложили и они еще не были разобраны:\n"))

		var problem *models.Problem
		var problemUrl string

		for i := range userSuggestions {
			problem = userSuggestions[i]
			problemUrl = problem.GetOriginalUrl()
			entityMessages = append(entityMessages, tu.Entity("\n- "))
			if problemUrl != "" {
				entityMessages = append(entityMessages, tu.Entity(string(problem.Name)).TextLink(problemUrl))
			} else {
				entityMessages = append(entityMessages, tu.Entity(string(problem.Name)))
			}
			entityMessages = append(entityMessages, tu.Entity("\n"))
		}
		botMessage = tu.MessageWithEntities(tu.ID(update.Message.Chat.ID), entityMessages...)
	} else {
		message := `
		Вы не предложили никакой задачи или все предложенные задачи уже разобраны.
		Самое время предложить что-то новое.
		/suggest_problem`
		botMessage = tu.Message(tu.ID(update.Message.Chat.ID), message)
	}

	_, _ = bot.SendMessage(botMessage)
}

// Обработчик команды "/suggest_problem". Формирует сообщение с просьюой предоставить ссылку на Problem
func ProcessProposeProblemFromMessage(bot *telego.Bot, update telego.Update) {
	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		string(custompredicates.GetMeProblemLink),
	).WithReplyMarkup(tu.ForceReply())

	_, _ = bot.SendMessage(message)
}

// Обработчик ошибок и отправка стандартного сообщения об ошибке.
// Позволяет добавлять кастомные ошибки models.ModelError.
func sendErrorMessage(bot *telego.Bot, update *telego.Update, err customerrors.CustomError) {
	logger := log.GetLogger()
	logger.Error(err)
	msg, ok := customerrors.CustomErrors[err]

	if !ok {
		msg = "Почему-то не получилось посмотреть информацию... Попробуйте позже."
	}

	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		msg,
	)
	_, _ = bot.SendMessage(message)
}
