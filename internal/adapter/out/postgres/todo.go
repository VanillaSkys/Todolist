package postgres

import (
	"github.com/VanillaSkys/todo_fiber/internal/core/domain/dto"
	"github.com/VanillaSkys/todo_fiber/internal/core/port/repository"
	"gorm.io/gorm"
)

type gormTodoRepositoryImpl struct {
	db *gorm.DB
}

func NewGormTodoRepository(db *gorm.DB) repository.TodoRepository {
	return &gormTodoRepositoryImpl{db: db}
}

func (g *gormTodoRepositoryImpl) FindAll() ([]dto.Todo, error) {
	var todos []dto.Todo
	result := g.db.Find(&todos)
	if result.Error != nil {
		return nil, result.Error
	}
	return todos, nil
}

func (g *gormTodoRepositoryImpl) Save(input dto.Todo) error {
	todo := dto.Todo{
		Id:          input.Id,
		Topic:       input.Topic,
		Description: input.Description,
		Status:      input.Status,
	}
	if result := g.db.Create(&todo); result.Error != nil {
		return result.Error
	}

	return nil
}
func (g *gormTodoRepositoryImpl) Update(input dto.TodoInputUpdateStatus) error {
	return nil
}

func (g *gormTodoRepositoryImpl) Delete(input dto.TodoInputDelete) error {
	return nil
}
