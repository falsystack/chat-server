package schema

import "time"

type Chat struct {
	ID      int       `json:"id"`
	Room    string    `json:"room"`
	Name    string    `json:"name"`
	Message string    `json:"message"`
	When    time.Time `json:"when"`
}
