package storages

import (
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
	"github.com/google/uuid"
)

const (
	selectSolution = `
SELECT id, name, problem_id, is_solved
FROM solutions
WHERE problem_id = $1
`

	insertSolvedSolution = `
BEGIN;
INSERT INTO solutions (name, problem_id, is_solved) VALUES ($1, $2, $3);
DELETE FROM suggestions WHERE problem_id = $2;
UPDATE problems SET status = $4 WHERE id = $2;
COMMIT;
`
	insertNotSolvedSolution = `
INSERT INTO solutions (name, problem_id, is_solved) VALUES ($1, $2, $3);
`

	updateSolutionStatusIsSolved = `
BEGGIN;
UPDATE solutions SET is_solved = true WHERE id = $1;
DELETE FROM suggestions WHERE problem_id = $2;
UPDATE problems SET status = $3 WHERE id = $2;
COMMIT;
`
)

func (s *Store) SolutionGetByProblemId(problemID uuid.UUID) (*models.Solution, error) {
	row := s.conn.QueryRow(s.ctx, selectSolution, problemID)

	var solution *models.Solution

	err := solution.ScanProblemRow(row)

	return solution, err
}

func (s *Store) CreateSolution(solution *models.Solution) error {
	var err error
	if solution.IsSolved {
		_, err = s.conn.Exec(s.ctx, insertSolvedSolution, solution.Name, solution.ProblemID, solution.IsSolved, models.CloseStatus)
	} else {
		_, err = s.conn.Exec(s.ctx, insertNotSolvedSolution, solution.Name, solution.ProblemID, solution.IsSolved, models.CloseStatus)
	}

	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	return nil
}

func (s *Store) MarkAsSolved(solution *models.Solution) error {
	_, err := s.conn.Exec(s.ctx, updateSolutionStatusIsSolved, solution.ID, solution.ProblemID, models.CloseStatus)

	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	return nil
}
