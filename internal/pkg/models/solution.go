package models

import (
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/customerrors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var solutionRepositoryURLs = map[string]string{
	"URLProblemDescription":  "https://github.com/DushnoAndTochka/solutions_algorithmic_problems/blob/main/solutions/%s/",
	"URLSolutionDescription": "https://github.com/DushnoAndTochka/solutions_algorithmic_problems/blob/main/solutions/%s/solution/",
}

// Храним решения проблем и название директории в которой лежит разбор проблемы.
type Solution struct {
	ID        uuid.UUID
	Name      string
	ProblemID uuid.UUID
	IsSolved  bool
}

func (s *Solution) GetURLProblemDESC() string {
	return fmt.Sprintf(solutionRepositoryURLs["URLProblemDescription"], s.Name)
}

func (s *Solution) GetURLSolutionDESC() (string, error) {
	if s.IsSolved {
		return fmt.Sprintf(solutionRepositoryURLs["URLSolutionDescription"], s.Name), nil
	} else {
		return "", customerrors.ErrSolutionNotReadyYet
	}
}

func (s *Solution) ScanRow(row pgx.Row) error {
	return row.Scan(
		&s.ID,
		&s.Name,
		&s.ProblemID,
		&s.IsSolved,
	)
}

func (s *Solution) ScanRows(rows pgx.Rows) error {
	return rows.Scan(
		&s.ID,
		&s.Name,
		&s.ProblemID,
		&s.IsSolved,
	)
}
