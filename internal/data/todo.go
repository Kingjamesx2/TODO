//filename: internal/data/todo.go

package data

import (
	"time"
)

type Todo struct {
	ID        int64     `json: "id"`
	CreatedAt time.Time `json: "created_at"`
	Name      string    `json: "name"`
	Task      string    `json: "task"`
}
