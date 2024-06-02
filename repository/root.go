package repository

import (
	"chat-server/config"
	"chat-server/types/schema"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type Repository struct {
	cfg *config.Config
	db  *sql.DB
}

const (
	room       = "chatting.room"
	chat       = "chatting.chat"
	serverInfo = "chatting.serverInfo"
)

func NewRepository(cfg *config.Config) (*Repository, error) {
	r := &Repository{cfg: cfg}
	var err error

	if _, err = sql.Open(cfg.DB.Database, cfg.DB.URL); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (r *Repository) Room(name string) (*schema.Room, error) {
	d := new(schema.Room)
	// select * from chatting.room where name = ?
	qs := query([]string{"SELECT * FROM", room, "WHERE name = >"})

	err := r.db.QueryRow(qs, name).Scan(&d)

	return d, err
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}
