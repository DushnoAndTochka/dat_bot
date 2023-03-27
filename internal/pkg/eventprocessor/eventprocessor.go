package eventprocessor

import (
	"regexp"

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
		"На данный момент я умею не так много. \n - Через меня можно предложить задачу на разбор\n - Посмотреть топ 5 желаемых задач от подписчиков.",
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
	answer := update.Message.Text
	var leetcodeValidLink = regexp.MustCompile(`^https:\/\/leetcode.com\/problems\/[a-z\-]*\/?$`)

	var message *telego.SendMessageParams
	if leetcodeValidLink.MatchString(answer) {
		message = tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Спасибо",
		)
	} else {
		message = tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Представленный вами URL не принадлежит LeetCode или указывает не на проблему.\nПример корректной ссылки: https://leetcode.com/problems/two-sum.",
		)
	}

	_, _ = bot.SendMessage(message)
}

func ProcessShowAllProposeProblems(bot *telego.Bot, update telego.Update) {
	// message := tu.Message(
	// 	tu.ID(update.Message.Chat.ID),
	// 	"Текущий Топ предложенных задач, выглядит следуюзим образом:\n",
	// ).WithReplyMarkup(tu.ForceReply())

	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"На этой неделе еще никто не успел предложить задачу. Вы можете быть первым.",
	)

	_, _ = bot.SendMessage(message)
}

func ProcessShowMyProposeProblem(bot *telego.Bot, update telego.Update) {
	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"На этой неделе вы не предложили никакой задачи. Самое время это сделать.\n/propose_problem",
	)

	_, _ = bot.SendMessage(message)
}

func ProcessProposeProblemFromMessage(bot *telego.Bot, update telego.Update) {
	message := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Укажите ссылку на задачу.",
	).WithReplyMarkup(tu.ForceReply())

	_, _ = bot.SendMessage(message)
}
