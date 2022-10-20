//filename: internal/data/todo.go

package data

import (
	"time"

	"todo.jamesfaber.net/internal/validator"
)

type Todo struct {
	ID        int64     `json: "id"`
	CreatedAt time.Time `json: "created_at"`
	Name      string    `json: "name"`
	Task      string    `json: "task"`
}

func ValidateTodo(v *validator.Validator, todo *Todo) {
	// Use the check() method to execute our validation checks
	v.Check(todo.Name != "", "name", "must be provided")
	v.Check(len(todo.Name) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(todo.Task != "", "task", "must be provided")
	v.Check(len(todo.Task) <= 200, "task", "must not be more than 200 bytes long")
}
