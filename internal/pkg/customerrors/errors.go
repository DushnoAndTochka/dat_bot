// Пакет определяющий кастомные ошибки. Содержит саму ошбку,
// которая будет отображаться в логах и текст, который будет отправлен пользователю.
package customerrors

import "fmt"

type CustomError error

// перечень кастомных ошибок
var (
	ErrNotSupportedURL     CustomError = fmt.Errorf("notSupported")
	ErrNotUniqSuggestion   CustomError = fmt.Errorf("notUniqSuggestion")
	ErrSolutionNotReadyYet CustomError = fmt.Errorf("solutionNotReadyYet")
)

// перечень сообщений, которые будут отправлены пользователю.
const (
	ErrNotSupportedURLMessage     string = "Данный URL не поддерживается"
	ErrNotUniqSuggestionMessage   string = "Вы уже предлагали данную задачу."
	ErrSolutionNotReadyYetMessage string = "Задача еще не имеет решения. Проверьте позднее или дождитесь публикации в канале @dushno_and_tochka."
)

var CustomErrors = map[CustomError]string{
	ErrNotSupportedURL:     ErrNotSupportedURLMessage,
	ErrNotUniqSuggestion:   ErrNotUniqSuggestionMessage,
	ErrSolutionNotReadyYet: ErrSolutionNotReadyYetMessage,
}
