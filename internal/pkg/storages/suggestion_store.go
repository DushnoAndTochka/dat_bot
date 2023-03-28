package storages

import (
	"context"
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

func (s *Store) SuggestionGetByTgID(ctx context.Context, suggestion *models.Suggestion) (*models.Suggestion, error) {
	row := s.conn.QueryRow(ctx, selectSuggestion, suggestion.UserID, suggestion.ProblemID)

	err := suggestion.ScanRow(row)

	return suggestion, err
}

func (s *Store) SuggestionCreate(ctx context.Context, suggestion *models.Suggestion) error {
	_, err := s.conn.Exec(ctx, insertSuggestion, suggestion.UserID, suggestion.ProblemID)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}
	return nil
}

func (s *Store) SuggestionUpdate(ctx context.Context, suggestion *models.Suggestion) error {
	_, err := s.conn.Exec(ctx, updateSuggestion, suggestion.UserID, suggestion.ProblemID)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}
	return nil
}
