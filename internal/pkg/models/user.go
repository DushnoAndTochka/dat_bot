package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mymmrac/telego"
)

const UserSelectValues = "id, name, tg_id"
const UserInsertValues = "name, tg_id"

type User struct {
	ID   uuid.UUID
	TgID int64
	Name string
}

func GetFromTg(update *telego.Update) *User {
	tgID := update.Message.From.ID
	name := update.Message.From.Username

	user := &User{
		TgID: tgID,
		Name: name,
	}

	return user
}

func (u *User) ScanUserRow(row pgx.Row) error {
	return row.Scan(
		&u.ID,
		&u.Name,
		&u.TgID,
	)
}

func (u *User) ScanUserRows(rows pgx.Rows) error {
	return rows.Scan(
		&u.ID,
		&u.Name,
		&u.TgID,
	)
}
