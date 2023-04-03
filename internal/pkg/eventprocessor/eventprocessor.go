package eventprocessor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/storages"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func ProcessStartComand(bot *telego.Bot, update telego.Update) {

	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Привет. Чем могут быть полезен ?",
	)

	_, _ = bot.SendMessage(message)
}

func ProcessHelpComand(bot *telego.Bot, update telego.Update) {

	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		`На данный момент я умею не так много. \n
		 - Через меня можно предложить задачу на разбор\n
		 - Посмотреть топ 5 желаемых задач от подписчиков.`,
	)

	_, _ = bot.SendMessage(message)
}

func ProcessNotSupportedComandsComand(bot *telego.Bot, update telego.Update) {
	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Неизвестная команда, используй /help.",
	).WithReplyMarkup(tu.ReplyKeyboardRemove())

	_, _ = bot.SendMessage(message)
}

func ProcessGetLinkFromReply(bot *telego.Bot, update telego.Update) {
	logger := log.GetLogger()
	answer := update.Message.Text
	user := models.GetFromTg(&update)
	storage := storages.GetStorage()
	err := storage.UserGetByTgID(user)

	if errors.Is(err, pgx.ErrNoRows) {
		if err = storage.UserCreate(user); err != nil {
			logger.Error(err)
			sendErrorMessage(bot, &update, err)

			return
		}
		user = models.GetFromTg(&update)
	} else if err != nil {
		logger.Errorf("store.UserGetByID: %v", err)
		sendErrorMessage(bot, &update, err)

		return
	}

	problem, err := models.NewProblemFromUrl(answer)

	var message *telego.SendMessageParams

	switch {
	case errors.Is(models.NotSupportedURL.Err, err):
		message = tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Представленный вами URL не принадлежит LeetCode или указывает не на проблему.\nПример корректной ссылки: https://leetcode.com/problems/two-sum.",
		)
		_, _ = bot.SendMessage(message)
		return
	}

	if err = storage.ProblemGet(problem); err != nil {
		logger.Error(err)
		sendErrorMessage(bot, &update, err)

		return
	}

	if problem.ID == uuid.Nil {
		if err = storage.ProblemCreate(problem); err != nil {
			logger.Error(err)
			sendErrorMessage(bot, &update, err)

			return
		}
		if err = storage.ProblemGet(problem); err != nil {
			logger.Error(err)
			sendErrorMessage(bot, &update, err)

			return
		}
	}

	suggection := models.NewSuggestion(user, problem)

	if err = storage.SuggestionGet(suggection); err != nil {
		logger.Error(err)
		sendErrorMessage(bot, &update, err)

		return
	}

	if suggection.ID != uuid.Nil {
		message = tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Вы уже предлагали данную задачу.",
		)
	} else {
		if err = storage.SuggestionCreate(suggection); err != nil {
			logger.Error(err)
			sendErrorMessage(bot, &update, err)

			return
		}
		message = tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Спасибо",
		)
	}

	_, _ = bot.SendMessage(message)
}

func ProcessShowAllProposeProblems(bot *telego.Bot, update telego.Update) {

	var message *telego.SendMessageParams

	logger := log.GetLogger()
	storage := storages.GetStorage()

	suggstions, err := storage.GetTopSuggestions()

	if err != nil {
		logger.Error("ProcessShowAllProposeProblems: %w", err)
		sendErrorMessage(bot, &update, err)
		return
	}

	if suggstions != nil {
		answerString := make([]string, len(suggstions)+1)
		answerString = append(answerString, "ТОП предложениями являются:\n")
		for k, v := range suggstions {
			answerString = append(answerString, fmt.Sprintf("Задача: '%s' была предложена %v раз.\n", string(*k), int(*v)))
		}
		message = tu.Message(
			tu.ID(update.Message.Chat.ID),
			strings.Join(answerString, "\n"),
		)
	} else {
		message = tu.Message(
			tu.ID(update.Message.Chat.ID),
			"На этой неделе еще никто не успел предложить задачу. Вы можете быть первым.",
		)
	}

	_, _ = bot.SendMessage(message)
}

func ProcessShowMyProposeProblem(bot *telego.Bot, update telego.Update) {
	var message *telego.SendMessageParams

	logger := log.GetLogger()
	user := models.GetFromTg(&update)
	storage := storages.GetStorage()
	err := storage.UserGetByTgID(user)

	if err != nil {
		logger.Errorf("store.UserGetByID: %v", err)
		sendErrorMessage(bot, &update, err)
		return
	}

	if user.ID == uuid.Nil {
		if err = storage.UserCreate(user); err != nil {
			logger.Error(err)
			return
		}
	}

	userSuggestions, err := storage.GetUserSuggestion(user)

	if err != nil {
		logger.Error("ProcessShowMyProposeProblem: %w", err)
		sendErrorMessage(bot, &update, err)
		return
	}

	if userSuggestions != nil {
		answerString := make([]string, len(userSuggestions)+1)
		answerString = append(answerString, "Задачи, которые вы предложили и они еще не были разобраны:\n")
		var problemName *models.ProblemName

		for i := range userSuggestions {
			problemName = userSuggestions[i]

			answerString = append(answerString, string(*problemName))
		}
		message = tu.Message(
			tu.ID(update.Message.Chat.ID),
			strings.Join(answerString, "\n"),
		)
	} else {
		message = tu.Message(
			tu.ID(update.Message.Chat.ID),
			`Вы не предложили никакой задачи или все предложенные задачи уже разобраны.\n
			Самое время предложить что-то новое.\n/propose_problem`,
		)
	}

	_, _ = bot.SendMessage(message)
}

func ProcessProposeProblemFromMessage(bot *telego.Bot, update telego.Update) {
	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Укажите ссылку на задачу.",
	).WithReplyMarkup(tu.ForceReply())

	_, _ = bot.SendMessage(message)
}

func sendErrorMessage(bot *telego.Bot, update *telego.Update, err error) {
	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Почему-то не получилось посмотреть информацию... Попробуйте позже.",
	)
	_, _ = bot.SendMessage(message)
}
