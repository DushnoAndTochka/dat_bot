package storages

import (
	"context"
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
)

var selectProblem = `
SELECT id, name, source
FROM problems
WHERE name = $1 and source = $2;
`

var insertProblem = `
INSERT INTO problems (name, source) VALUES ($1, $2);
`

func (s *Store) ProblemGetByTgID(ctx context.Context, problem *models.Problem) (*models.Problem, error) {
	row := s.conn.QueryRow(ctx, selectProblem, problem.Name, problem.Source)

	err := problem.ScanProblemRow(row)

	return problem, err
}

func (s *Store) ProblemCreate(ctx context.Context, problem *models.Problem) error {
	_, err := s.conn.Exec(ctx, insertProblem, problem.Name, problem.Source)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}
	return nil
}
