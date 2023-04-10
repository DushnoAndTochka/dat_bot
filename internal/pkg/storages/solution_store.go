package storages

import (
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
	"github.com/google/uuid"
)

const (
	selectSolution = `
SELECT id, name, problem_id, is_solved
FROM solutions
WHERE problem_id = $1
`

	selectAllSolutions = `
SELECT id, name, problem_id, is_solved
FROM solutions
`

	insertSolution = `
INSERT INTO solutions (name, problem_id, is_solved) VALUES ($1, $2, $3);
`
	updateSolutionStatus = `
UPDATE solutions SET is_solved = $2 WHERE id = $1;
`
)

func (s *Store) SolutionGetByProblemId(problemID uuid.UUID) (*models.Solution, error) {

	logger := log.GetLogger()
	row := s.conn.QueryRow(s.ctx, selectSolution, problemID)

	solution := &models.Solution{}

	err := solution.ScanRow(row)
	if err != nil {
		logger.Error("SolutionGetByProblemId: %w", err)
	}

	return solution, err
}

func (s *Store) SolutionsGetAll() (map[string]*models.Solution, error) {
	logger := log.GetLogger()
	rows, err := s.conn.Query(s.ctx, selectAllSolutions)
	if err != nil {
		logger.Error("SolutionStore SolutionsGetAll failed: %w", err)
		return nil, err
	}

	solutions := make(map[string]*models.Solution)
	var solution *models.Solution

	for rows.Next() {
		solution = &models.Solution{}
		err = solution.ScanRows(rows)
		if err != nil {
			logger.Error("SolutionStore SolutionsGetAll ScanRow failed: %w", err)
		}
		solutions[solution.Name] = solution
	}

	return solutions, nil
}

func (s *Store) SolutionUpdateOrCreate(solution *models.Solution) error {
	logger := log.GetLogger()
	tx, err := s.conn.Begin(s.ctx)

	if err != nil {
		tx.Rollback(s.ctx)
		logger.Error(err)
		return err
	}

	if solution.ID == uuid.Nil {
		_, err = tx.Exec(s.ctx, insertSolution, solution.Name, solution.ProblemID, solution.IsSolved)
		if err != nil {
			logger.Error("SolutionStore CreateOrUpdate CreateSolution is failed: %w", err)
			tx.Rollback(s.ctx)
			return err
		}
	} else {
		_, err = tx.Exec(s.ctx, updateSolutionStatus, solution.ID, solution.IsSolved)
		if err != nil {
			logger.Error("SolutionStore CreateOrUpdate UpdateSolution is failed: %w", err)
			tx.Rollback(s.ctx)
			return err
		}
	}

	_, err = tx.Exec(s.ctx, deleteSuggestionByProblem, solution.ProblemID)
	if err != nil {
		logger.Error("SolutionStore CreateOrUpdate DeleteSuggestion is failed: %w", err)
		tx.Rollback(s.ctx)
		return err
	}

	_, err = tx.Exec(s.ctx, updateProblemStatus, solution.ProblemID, models.CloseStatus)
	if err != nil {
		logger.Error("SolutionStore CreateOrUpdate UpdateProblemStatus is failed: %w", err)
		tx.Rollback(s.ctx)
		return err
	}
	tx.Commit(s.ctx)
	return nil
}
