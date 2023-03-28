package storages

import (
	"context"
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

func (s *Store) UserGetByTgID(ctx context.Context, tgID int64) (*models.User, error) {
	row := s.conn.QueryRow(ctx, selectUser, tgID)
	u := &models.User{}
	err := u.ScanUserRow(row)

	return u, err
}

func (s *Store) UserCreate(ctx context.Context, u *models.User) error {
	_, err := s.conn.Exec(ctx, insertUser, u.TgID, u.Name)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}
	return nil
}
