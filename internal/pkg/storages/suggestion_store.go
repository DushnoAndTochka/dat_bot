package storages

import (
	"errors"
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/customerrors"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	// поиск предложения по user_id и problem_id
	selectSuggestion = `
SELECT id, user_id, problem_id
FROM suggestions
WHERE user_id = $1 and problem_id = $2;
`

	// поиск всех предложений пользователя
	selectUserSuggestions = `
SELECT suggestions.problem_id, 
       problems.name,
	   problems.source,
	   problems.status
FROM suggestions
JOIN problems ON suggestions.problem_id = problems.id
WHERE user_id = $1;
`

	// поиск топ 10 предложений
	selectTOPSuggestions = `
SELECT suggestions.problem_id, 
       count(suggestions.problem_id) as countSuggestions, 
       problems.name,
	   problems.source
FROM suggestions
JOIN problems ON suggestions.problem_id = problems.id
WHERE problems.status = $1
GROUP BY suggestions.problem_id, problems.name, problems.source
ORDER BY countSuggestions DESC
LIMIT 10;
`

	// Создание нового предложения
	insertSuggestion = `
INSERT INTO suggestions (user_id, problem_id) VALUES ($1, $2);
`

	// Обновление предложения
	updateSuggestion = `
UPDATE suggestions
SET problem_id = $2
WHERE id = $1;
`
)

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

// Проверяет наличие Suggestion. Если он есть, то ничего не делает, иначе создает новое предложение.
func (s *Store) SuggestionCheckOrCreate(suggestion *models.Suggestion) error {
	logger := log.GetLogger()
	err := s.SuggestionGet(suggestion)
	if errors.Is(err, pgx.ErrNoRows) {
		logger.Info("SuggestionCheckOrCreate: Try to create new suggestion.")
		if err := s.SuggestionCreate(suggestion); err != nil {
			logger.Error("SuggestionCheckOrCreate: Create failed: %w.", err)
			return err
		}
		logger.Info("SuggestionCheckOrCreate: SuggestionCreate is succeeded.")
		err = s.SuggestionGet(suggestion)
	} else if suggestion.ID != uuid.Nil {
		logger.Info("SuggestionCheckOrCreate: Is not uniq suggestion.")
		return customerrors.ErrNotUniqSuggestion
	}

	return err
}

func (s *Store) SuggestionUpdate(suggestion *models.Suggestion) error {
	_, err := s.conn.Exec(s.ctx, updateSuggestion, suggestion.UserID, suggestion.ProblemID)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}
	return nil
}

func (s *Store) GetTopSuggestions() ([]*models.Problem, error) {
	logger := log.GetLogger()
	problemStatus := models.OpenStatus
	rows, err := s.conn.Query(s.ctx, selectTOPSuggestions, problemStatus)

	if err != nil {
		return nil, err
	}

	var topSuggestions []*models.Problem

	var problemID *uuid.UUID
	var problemName string
	var problemSource string
	var problem *models.Problem
	var countSuggestions *models.CountSuggestions

	for rows.Next() {

		err = rows.Scan(&problemID, &countSuggestions, &problemName, &problemSource)
		logger.Debug("GetTopSuggestions: Result: ", countSuggestions, " ", *problemID, " ", problemName, " ", problemSource, " ", string(problemStatus))
		if err != nil {
			logger.Debug("GetTopSuggestions: ", err)
		}
		problem, err = models.NewProblem(*problemID, problemName, problemSource, string(problemStatus))
		problem.CountSuggestions = *countSuggestions

		if err != nil {
			logger.Debug("GetTopSuggestions: ", err)
		}
		topSuggestions = append(topSuggestions, problem)
	}

	return topSuggestions, err
}

func (s *Store) GetUserSuggestion(user *models.User) ([]*models.Problem, error) {

	rows, err := s.conn.Query(s.ctx, selectUserSuggestions, user.ID)

	var userSuggestions []*models.Problem

	if err != nil {
		return userSuggestions, err
	}

	var problem *models.Problem

	for rows.Next() {
		problem = &models.Problem{}
		err = problem.ScanProblemRows(rows)

		userSuggestions = append(userSuggestions, problem)
	}

	return userSuggestions, err
}
