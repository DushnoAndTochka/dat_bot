package storages

import (
	"errors"
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var selectProblem = `
SELECT id, name, source, status
FROM problems
WHERE name = $1 and source = $2;
`

var selectAllProblemsWithCurrentStatus = `
SELECT id FROM problems
WHERE status = $1
`

var insertProblem = `
INSERT INTO problems (name, source, status) VALUES ($1, $2, $3);
`

func (s *Store) ProblemGet(problem *models.Problem) error {
	logger := log.GetLogger()
	row := s.conn.QueryRow(s.ctx, selectProblem, problem.Name, problem.Source)

	err := problem.ScanProblemRow(row)

	if errors.Is(err, pgx.ErrNoRows) {
		logger.Info("store.ProblemGet: Not Found problem name: %s, source: %s", problem.Name, problem.Source)
		return nil
	}

	return err
}

func (s *Store) ProblemCreate(problem *models.Problem) error {
	_, err := s.conn.Exec(s.ctx, insertProblem, problem.Name, problem.Source, problem.Status)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}
	return nil
}

func (s *Store) ProblemGetAllIDSOpen() ([]*uuid.UUID, error) {

	var problemIDs []*uuid.UUID

	rows, err := s.conn.Query(s.ctx, selectAllProblemsWithCurrentStatus, models.OpenStatus)

	if err != nil {
		return nil, fmt.Errorf("conn.Query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id *uuid.UUID
		err = rows.Scan(&id)
		problemIDs = append(problemIDs, (*uuid.UUID)(id))
	}

	return problemIDs, err
}
