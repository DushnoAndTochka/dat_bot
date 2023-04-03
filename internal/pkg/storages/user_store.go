package storages

import (
	"fmt"

	"github.com/artem-telnov/dushno_and_tochka_bot/internal/pkg/models"
)

var selectUser = `
SELECT id, name, tg_id
FROM users
WHERE tg_id = $1;
`

var insertUser = `
INSERT INTO users (name, tg_id) VALUES ($1, $2);
`

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
