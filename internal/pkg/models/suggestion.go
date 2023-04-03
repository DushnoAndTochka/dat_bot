package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const SuggestionSelectValues = "id, user_id, problem_id"
const SuggestionInsertValues = "user_id, problem_id"

type CountSuggestions int

type Suggestion struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ProblemID uuid.UUID
}

func NewSuggestion(user *User, problem *Problem) *Suggestion {
	return &Suggestion{
		UserID:    user.ID,
		ProblemID: problem.ID,
	}
}

func (s *Suggestion) ScanRow(row pgx.Row) error {
	return row.Scan(
		&s.ID,
		&s.UserID,
		&s.ProblemID,
	)
}

func (s *Suggestion) ScanRows(rows pgx.Rows) error {
	return rows.Scan(
		&s.ID,
		&s.UserID,
		&s.ProblemID,
	)
}
