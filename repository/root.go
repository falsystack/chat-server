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

func (r Repository) GetChatList(roomName string) ([]*schema.Chat, error) {
	qs := query([]string{"SELECT * FROM", chat, "WHERE room = ? ORDER BY `when` DESC LIMIT 10"})
	if cursor, err := r.db.Query(qs, roomName); err != nil {
		return nil, err
	} else {
		defer cursor.Close()

		var result []*schema.Room

		for cursor.Next() {
			d := new(schema.Chat)

			if err := cursor.Scan(&d); err != nil {
				return nil, err
			} else {
				result = append(result, d)
			}
		}

		if len(result) == 0 {
			return []*schema.Room{}, nil
		} else {
			return result, nil
		}
	}
}

func (r *Repository) RoomList() ([]*schema.Room, error) {
	qs := query([]string{"SELECT * FROM", room})
	if cursor, err := r.db.Query(qs); err != nil {
		return nil, err
	} else {
		defer cursor.Close()

		var result []*schema.Room

		for cursor.Next() {
			d := new(schema.Chat)

			if err := cursor.Scan(&d); err != nil {
				return nil, err
			} else {
				result = append(result, d)
			}
		}

		if len(result) == 0 {
			return []*schema.Room{}, nil
		} else {
			return result, nil
		}
	}
}

func (r Repository) MakeRoom(name string) error {
	_, err := r.db.Exec("INSERT INTO chatting.room(name) VALUES(?)", name)
	return err
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
