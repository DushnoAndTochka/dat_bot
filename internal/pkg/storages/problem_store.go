package storages

import (
	"errors"
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
	"github.com/jackc/pgx/v5"
)

const (
	// поиск проблем по имени и ресурсу.
	selectProblem = `
SELECT id, name, source, status
FROM problems
WHERE name = $1 and source = $2;
`

	// Создание новой проблемы
	insertProblem = `
INSERT INTO problems (name, source, status) VALUES ($1, $2, $3);
`

	// Обновление статуса проблемы
	updateProblemStatus = `
UPDATE problems SET status = $2 WHERE id = $1;
`
)

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

// Метод пытается найти объект в бд. Если не получилось, то создает новый объект.
func (s *Store) ProblemGetOrCreate(problem *models.Problem) error {
	logger := log.GetLogger()
	err := s.ProblemGet(problem)
	if errors.Is(err, pgx.ErrNoRows) {
		logger.Info("ProblemGetOrCreate: Problem is not found. Try to create new Problem")
		if err := s.ProblemCreate(problem); err != nil {
			logger.Error("ProblemGetOrCreate: CreateProblem is failed: %w", err)
			return err
		}
		logger.Info("ProblemGetOrCreate: CreateProblem is succeeded.")
		err = s.ProblemGet(problem)
	}

	return err
}
