package storages

import (
	"errors"
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/log"
	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	selectUser = `
SELECT id, name, tg_id
FROM users
WHERE tg_id = $1;
`

	insertUser = `
INSERT INTO users (name, tg_id) VALUES ($1, $2);
`
)

func (s *Store) UserGetByTgID(u *models.User) error {
	row := s.conn.QueryRow(s.ctx, selectUser, u.TgID)
	err := u.ScanUserRow(row)

	return err
}

func (s *Store) UserCreate(u *models.User) error {
	_, err := s.conn.Exec(s.ctx, insertUser, u.Name, u.TgID)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}
	return nil
}

// Ищет или создает.
func (s *Store) UserGetOrCreate(u *models.User) error {
	logger := log.GetLogger()
	err := s.UserGetByTgID(u)

	if errors.Is(err, pgx.ErrNoRows) || u.ID == uuid.Nil {
		logger.Info("UserGetOrCreate: User not gound. Try to create new user.")
		if err = s.UserCreate(u); err != nil {
			logger.Error("UserGetOrCreate: Create user failed: %w", err)
			return err
		}
		logger.Info("UserGetOrCreate: UserCreate is succeeded.")
		err = s.UserGetByTgID(u)
	}
	logger.Info("UserGetOrCreate: User successfuly found. %w", u)

	return err
}
