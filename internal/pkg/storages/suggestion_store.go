package storages

import (
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
)

var selectSuggestion = `
SELECT id, user_id, problem_id
FROM suggestions
WHERE name = $1 and source = $2;
`

var insertSuggestion = `
INSERT INTO suggestions (user_id, problem_id) VALUES ($1, $2);
`

var updateSuggestion = `
UPDATE suggestions
SET problem_id = $2
WHERE id = $1;
`

func (s *Store) SuggestionGet(suggestion *models.Suggestion) error {
	row := s.conn.QueryRow(s.ctx, selectSuggestion, suggestion.UserID, suggestion.ProblemID)

	err := suggestion.ScanRow(row)

	return err
}

func (s *Store) SuggestionCreate(suggestion *models.Suggestion) error {
	_, err := s.conn.Exec(s.ctx, insertSuggestion, suggestion.UserID, suggestion.ProblemID)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}
	return nil
}

func (s *Store) SuggestionUpdate(suggestion *models.Suggestion) error {
	_, err := s.conn.Exec(s.ctx, updateSuggestion, suggestion.UserID, suggestion.ProblemID)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}
	return nil
}
