package custompredicates

import "github.com/mymmrac/telego"

// Сообщения которые будут использоваться в кастомных условиях для Handlers.
type CustomPredicateMessage string

const (
	GetMeProblemLink CustomPredicateMessage = "Укажите ссылку на задачу."
)

// кастомное правило для обработчиков бота
// позволяет определить, что бот написал сообщение "Укажите ссылку на задачу.",
// а пользователь на него ответил;
func IsNewProposeTask(update telego.Update) bool {
	if update.Message != nil &&
		update.Message.ReplyToMessage != nil &&
		update.Message.ReplyToMessage.From.IsBot {
		return update.Message.ReplyToMessage.Text == string(GetMeProblemLink)
	}
	return false
}
