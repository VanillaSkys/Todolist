package repository

import "github.com/VanillaSkys/todo_fiber/internal/core/domain/dto"

type TodoRepository interface {
	FindAll() ([]dto.Todo, error)
	Save(dto.Todo) error
	Update(dto.TodoInputUpdateStatus) error
	Delete(dto.TodoInputDelete) error
}
