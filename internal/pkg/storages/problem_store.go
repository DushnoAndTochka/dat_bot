package storages

import (
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
)

var selectProblem = `
SELECT id, name, source, status
FROM problems
WHERE name = $1 and source = $2;
`

var insertProblem = `
INSERT INTO problems (name, source, status) VALUES ($1, $2, $3);
`

func (s *Store) ProblemGet(problem *models.Problem) error {
	row := s.conn.QueryRow(s.ctx, selectProblem, problem.Name, problem.Source)

	err := problem.ScanProblemRow(row)

	return err
}

func (s *Store) ProblemCreate(problem *models.Problem) error {
	_, err := s.conn.Exec(s.ctx, insertProblem, problem.Name, problem.Source, problem.Status)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}
	return nil
}
