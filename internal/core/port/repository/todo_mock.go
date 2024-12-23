package repository

import (
	"github.com/VanillaSkys/todo_fiber/internal/core/domain/dto"
	"github.com/stretchr/testify/mock"
)

type todoRepositoryMock struct {
	mock.Mock
}

func NewTodoRepositoryMock() *todoRepositoryMock {
	return &todoRepositoryMock{}
}

func (m *todoRepositoryMock) FindAll() ([]dto.Todo, error) {
	args := m.Called()
	return args.Get(0).([]dto.Todo), args.Error(1)
}

func (m *todoRepositoryMock) Save(input dto.Todo) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *todoRepositoryMock) Update(input dto.TodoInputUpdateStatus) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *todoRepositoryMock) Delete(input dto.TodoInputDelete) error {
	args := m.Called(input)
	return args.Error(0)
}
