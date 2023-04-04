package models

import "fmt"

type ModelError error

var (
	ErrNotSupportedURL   ModelError = fmt.Errorf("notSupported")
	ErrNotUniqSuggestion ModelError = fmt.Errorf("notUniqSuggestion")
)

const (
	ErrNotSupportedURLMessage   string = "Данный URL не поддерживается"
	ErrNotUniqSuggestionMessage string = "Вы уже предлагали данную задачу."
)

var ModelErrors = map[ModelError]string{
	ErrNotSupportedURL:   ErrNotSupportedURLMessage,
	ErrNotUniqSuggestion: ErrNotUniqSuggestionMessage,
}
