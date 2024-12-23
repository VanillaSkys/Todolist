package service

import (
	"context"
	"encoding/json"

	"github.com/VanillaSkys/todo_fiber/internal/core/domain/dto"
	"github.com/VanillaSkys/todo_fiber/internal/core/port/cache"
	"github.com/VanillaSkys/todo_fiber/internal/core/port/repository"
)

type TodoService interface {
	FindAll() ([]dto.Todo, error)
	Create(dto.Todo) error
	Update(dto.TodoInputUpdateStatus) error
	Delete(dto.TodoInputDelete) error
}

type todoServiceImpl struct {
	repo  repository.TodoRepository
	cache cache.Cache
}

func NewTodoService(repo repository.TodoRepository, cache cache.Cache) TodoService {
	return &todoServiceImpl{
		repo:  repo,
		cache: cache,
	}
}

func (s *todoServiceImpl) FindAll() ([]dto.Todo, error) {
	cachedData, err := s.cache.Get(context.Background(), "todos")

	if err != nil {
		todos, err := s.repo.FindAll()
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(todos)
		if err != nil {
			return nil, err
		}

		if err := s.cache.Set(context.Background(), "todos", string(data), 0); err != nil {
			return nil, err
		}

		return todos, nil
	}

	var todos []dto.Todo
	if err := json.Unmarshal([]byte(cachedData), &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *todoServiceImpl) Create(input dto.Todo) error {
	if err := s.repo.Save(input); err != nil {
		return err
	}
	cachedData, err := s.cache.Get(context.Background(), "todos")
	var newData []dto.Todo

	if err != nil {
		// If cache miss, fetch all todos from DB
		todos, err := s.repo.FindAll()
		if err != nil {
			return err
		}

		newData = todos
	} else {
		// Cache hit, unmarshal and append the new todo
		if err := json.Unmarshal([]byte(cachedData), &newData); err != nil {
			return err
		}
		newData = append(newData, input)
	}

	// Marshal and update the cache
	data, err := json.Marshal(newData)
	if err != nil {
		return err
	}

	return s.cache.Set(context.Background(), "todos", string(data), 0)
}

func (s *todoServiceImpl) Update(input dto.TodoInputUpdateStatus) error {
	if err := s.repo.Update(input); err != nil {
		return err
	}
	cacheData, err := s.cache.Get(context.Background(), "todos")

	if err != nil {
		todos, err := s.repo.FindAll()

		if err != nil {
			return err
		}

		data, err := json.Marshal(todos)
		if err != nil {
			return err
		}
		return s.cache.Set(context.Background(), "todos", string(data), 0)
	}

	var todos []dto.Todo
	if err := json.Unmarshal([]byte(cacheData), &todos); err != nil {
		return err
	}

	for index, todo := range todos {
		if todo.Id == input.Id {
			todos[index].Status = input.Status
			break
		}
	}

	data, err := json.Marshal(todos)
	if err != nil {
		return err
	}

	return s.cache.Set(context.Background(), "todos", string(data), 0)
}

func (s *todoServiceImpl) Delete(input dto.TodoInputDelete) error {
	err := s.repo.Delete(input)
	if err != nil {
		return err
	}

	cacheData, err := s.cache.Get(context.Background(), "todos")
	if err != nil {
		todos, err := s.repo.FindAll()
		if err != nil {
			return err
		}

		data, err := json.Marshal(todos)
		if err != nil {
			return err
		}
		return s.cache.Set(context.Background(), "todos", string(data), 0)
	}

	var todos []dto.Todo

	if err := json.Unmarshal([]byte(cacheData), &todos); err != nil {
		return err
	}

	var newData []dto.Todo
	for _, todo := range todos {
		if todo.Id != input.Id {
			newData = append(newData, todo)
		}
	}

	data, err := json.Marshal(newData)
	if err != nil {
		return err
	}

	return s.cache.Set(context.Background(), "todos", string(data), 0)
}
