package storages

import (
	"errors"
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var selectSuggestion = `
SELECT id, user_id, problem_id
FROM suggestions
WHERE user_id = $1 and problem_id = $2;
`

var selectUserSuggestions = `
SELECT suggestions.problem_id, 
       problems.name
FROM suggestions
JOIN problems ON suggestions.problem_id = problems.id
WHERE user_id = $1;
`

var selectTOPSuggestions = `
SELECT suggestions.problem_id, 
       count(suggestions.problem_id) as countSuggestions, 
       problems.name
FROM suggestions
JOIN problems ON suggestions.problem_id = problems.id
WHERE problems.status = $1
GROUP BY suggestions.problem_id, problems.name
ORDER BY countSuggestions
LIMIT 10;
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

	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}

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

func (s *Store) GetTopSuggestions() (map[*models.ProblemName]*models.CountSuggestions, error) {
	rows, err := s.conn.Query(s.ctx, selectTOPSuggestions, models.OpenStatus)

	if err != nil {
		return nil, err
	}

	topSuggestions := make(map[*models.ProblemName]*models.CountSuggestions)

	for rows.Next() {
		var problemID *uuid.UUID
		var problemName *models.ProblemName
		var countSuggestions *models.CountSuggestions

		err = rows.Scan(&problemID, &countSuggestions, &problemName)

		topSuggestions[problemName] = countSuggestions
	}

	return topSuggestions, err
}

func (s *Store) GetUserSuggestion(user *models.User) ([]*models.ProblemName, error) {
	rows, err := s.conn.Query(s.ctx, selectUserSuggestions, user.ID)

	if err != nil {
		return nil, err
	}

	var userSuggestions []*models.ProblemName

	for rows.Next() {
		var problemID *uuid.UUID
		var problemName *models.ProblemName

		err = rows.Scan(&problemID, &problemName)

		userSuggestions = append(userSuggestions, problemName)
	}

	return userSuggestions, err
}
