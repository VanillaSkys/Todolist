package service_test

import (
	"errors"
	"testing"

	"github.com/VanillaSkys/todo_fiber/internal/core/domain/dto"
	"github.com/VanillaSkys/todo_fiber/internal/core/domain/service"
	"github.com/VanillaSkys/todo_fiber/internal/core/port/cache"
	"github.com/VanillaSkys/todo_fiber/internal/core/port/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTodoserviceFindAllByTodo(t *testing.T) {
	testCases := []struct {
		description string
		repoReturn  struct {
			todos []dto.Todo
			err   error
		}
		cacheGetReturn struct {
			data string
			err  error
		}
		cacheSetReturn error
		expected       []dto.Todo
		expectedErr    error
	}{
		// Cache hit scenario
		{
			description: "Cache hit",
			repoReturn: struct {
				todos []dto.Todo
				err   error
			}{
				todos: []dto.Todo{
					{
						Id:          "1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412",
						Topic:       "Complete Project",
						Description: "Description for Complete Project",
						Status:      "Completed",
					},
				},
				err: nil,
			},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "[{\"id\":\"1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412\",\"topic\":\"Complete Project\",\"description\":\"Description for Complete Project\",\"status\":\"Completed\"}]",
				err:  nil,
			},
			cacheSetReturn: nil,
			expected: []dto.Todo{
				{
					Id:          "1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412",
					Topic:       "Complete Project",
					Description: "Description for Complete Project",
					Status:      "Completed",
				},
			},
			expectedErr: nil,
		},
		// Cache miss scenario
		{
			description: "Cache miss",
			repoReturn: struct {
				todos []dto.Todo
				err   error
			}{
				todos: []dto.Todo{
					{
						Id:          "1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412",
						Topic:       "Complete Project",
						Description: "Description for Complete Project",
						Status:      "Completed",
					},
				},
				err: nil,
			},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "", // Simulate cache miss
				err:  errors.New("failed to get cache"),
			},
			cacheSetReturn: nil,
			expected: []dto.Todo{
				{
					Id:          "1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412",
					Topic:       "Complete Project",
					Description: "Description for Complete Project",
					Status:      "Completed",
				},
			},
			expectedErr: nil,
		},
		// Repository error scenario
		{
			description: "Repository error",
			repoReturn: struct {
				todos []dto.Todo
				err   error
			}{
				todos: nil,
				err:   errors.New("failed to fetch todos"),
			},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "",
				err:  errors.New("failed to get cache"),
			},
			cacheSetReturn: nil,
			expected:       nil,
			expectedErr:    errors.New("failed to fetch todos"),
		},
		// Cache error scenario
		{
			description: "Cache error",
			repoReturn: struct {
				todos []dto.Todo
				err   error
			}{
				todos: []dto.Todo{
					{
						Id:          "1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412",
						Topic:       "Complete Project",
						Description: "Description for Complete Project",
						Status:      "Completed",
					},
				},
				err: nil,
			},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "",
				err:  errors.New("failed to get cache"),
			},
			cacheSetReturn: errors.New("failed to set cache"),
			expected:       nil,
			expectedErr:    errors.New("failed to set cache"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			// Arrange
			todoRepo := repository.NewTodoRepositoryMock()

			todoCache := cache.NewRedisCacheMock()

			todoCache.On("Get", mock.Anything, "todos").Return(testCase.cacheGetReturn.data, testCase.cacheGetReturn.err)

			if testCase.cacheGetReturn.err != nil {
				todoRepo.On("FindAll").Return(testCase.repoReturn.todos, testCase.repoReturn.err)
				if testCase.repoReturn.err == nil {
					todoCache.On("Set", mock.Anything, "todos", mock.Anything, mock.Anything).Return(testCase.cacheSetReturn)
				}
			}

			todoService := service.NewTodoService(todoRepo, todoCache)

			// Act
			response, err := todoService.FindAll()

			// Assert
			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, response)
			}
			todoRepo.AssertExpectations(t)
			todoCache.AssertExpectations(t)
		})
	}
}

func TestTodoserviceCreateTodo(t *testing.T) {
	testCases := []struct {
		description       string
		input             dto.Todo
		repoSaveReturn    error
		repoFindAllReturn struct {
			todos []dto.Todo
			err   error
		}
		cacheGetReturn struct {
			data string
			err  error
		}
		cacheSetReturn error
		expectedErr    error
	}{
		{
			description: "Create todo success repo and cache hit",
			input: dto.Todo{
				Id:          "1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412",
				Topic:       "Complete Project",
				Description: "Description for Complete Project",
				Status:      "Completed",
			},
			repoSaveReturn: nil,
			repoFindAllReturn: struct {
				todos []dto.Todo
				err   error
			}{
				todos: []dto.Todo{
					{
						Id:          "1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412",
						Topic:       "Complete Project",
						Description: "Description for Complete Project",
						Status:      "Completed",
					},
				},
				err: nil,
			},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "[{\"id\":\"1e89f1d7-78c5-4d4a-bae3-d4f5f96a7411\",\"topic\":\"Complete Project\",\"description\":\"Description for Complete Project\",\"status\":\"Completed\"}]",
				err:  nil,
			},
			cacheSetReturn: nil,
			expectedErr:    nil,
		},
		{
			description: "Create todo failed due to repository failure",
			input: dto.Todo{
				Id:          "1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412",
				Topic:       "Complete Project",
				Description: "Description for Complete Project",
				Status:      "Completed",
			},
			repoSaveReturn: errors.New("repository save failed"),
			repoFindAllReturn: struct {
				todos []dto.Todo
				err   error
			}{
				todos: []dto.Todo{},
				err:   nil,
			},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "[{\"id\":\"1e89f1d7-78c5-4d4a-bae3-d4f5f96a7411\",\"topic\":\"Complete Project\",\"description\":\"Description for Complete Project\",\"status\":\"Completed\"}]",
				err:  nil,
			},
			cacheSetReturn: nil,
			expectedErr:    errors.New("repository save failed"),
		},
		{
			description: "Cache miss get",
			input: dto.Todo{
				Id:          "1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412",
				Topic:       "Complete Project",
				Description: "Description for Complete Project",
				Status:      "Completed",
			},
			repoSaveReturn: nil,
			repoFindAllReturn: struct {
				todos []dto.Todo
				err   error
			}{
				todos: []dto.Todo{},
				err:   nil,
			},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "",
				err:  errors.New("cache miss get"),
			},
			cacheSetReturn: nil,
			expectedErr:    nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			// Arrange
			todoRepo := repository.NewTodoRepositoryMock()
			todoCache := cache.NewRedisCacheMock()
			todoRepo.On("Save", testCase.input).Return(testCase.repoSaveReturn)
			if testCase.repoSaveReturn == nil {

				todoCache.On("Get", mock.Anything, "todos").Return(testCase.cacheGetReturn.data, testCase.cacheGetReturn.err)

				if testCase.cacheGetReturn.err != nil {
					todoRepo.On("FindAll").Return(testCase.repoFindAllReturn.todos, testCase.repoFindAllReturn.err)
				}

				todoCache.On("Set", mock.Anything, "todos", mock.Anything, mock.Anything).Return(testCase.cacheSetReturn)
			}

			todoService := service.NewTodoService(todoRepo, todoCache)

			// Act
			err := todoService.Create(testCase.input)

			// Assert
			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			todoRepo.AssertExpectations(t)
		})
	}
}

func TestTodoserviceUpdateStatus(t *testing.T) {
	testCases := []struct {
		description       string
		input             dto.TodoInputUpdateStatus
		repoUpdateReturn  error
		repoFindAllReturn struct {
			todos []dto.Todo
			err   error
		}
		cacheGetReturn struct {
			data string
			err  error
		}
		cacheSetReturn error
		expectedErr    error
	}{
		{
			description: "Update status success.",
			input: dto.TodoInputUpdateStatus{
				Id:     "1",
				Status: "Pending",
			},
			repoUpdateReturn: nil,
			repoFindAllReturn: struct {
				todos []dto.Todo
				err   error
			}{
				todos: []dto.Todo{
					{
						Id:          "1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412",
						Topic:       "Complete Project",
						Description: "Description for Complete Project",
						Status:      "Completed",
					},
				},
				err: nil,
			},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "[{\"id\":\"1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412\",\"topic\":\"Complete Project\",\"description\":\"Description for Complete Project\",\"status\":\"Completed\"}]",
				err:  nil,
			},
			cacheSetReturn: nil,
			expectedErr:    nil,
		},
		{
			description: "Update status is miss cache.",
			input: dto.TodoInputUpdateStatus{
				Id:     "1",
				Status: "Pending",
			},
			repoUpdateReturn: nil,
			repoFindAllReturn: struct {
				todos []dto.Todo
				err   error
			}{
				todos: []dto.Todo{
					{
						Id:          "1e89f1d7-78c5-4d4a-bae3-d4f5f96a7412",
						Topic:       "Complete Project",
						Description: "Description for Complete Project",
						Status:      "Completed",
					},
				},
				err: nil,
			},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "",
				err:  errors.New("miss cache"),
			},
			cacheSetReturn: nil,
			expectedErr:    nil,
		},
		{
			description: "Update status is failed repository update dont have Id.",
			input: dto.TodoInputUpdateStatus{
				Id:     "not",
				Status: "Complete",
			},
			repoUpdateReturn: errors.New("error update"),
			repoFindAllReturn: struct {
				todos []dto.Todo
				err   error
			}{
				todos: []dto.Todo{},
				err:   nil,
			},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "",
				err:  nil,
			},
			cacheSetReturn: nil,
			expectedErr:    errors.New("error update"),
		},
		{
			description: "Update status is failed repository findall.",
			input: dto.TodoInputUpdateStatus{
				Id:     "1",
				Status: "Pending",
			},
			repoUpdateReturn: nil,
			repoFindAllReturn: struct {
				todos []dto.Todo
				err   error
			}{
				todos: []dto.Todo{},
				err:   errors.New("error findall"),
			},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "",
				err:  errors.New("miss cache."),
			},
			cacheSetReturn: nil,
			expectedErr:    errors.New("error findall"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			// Arrange
			todoRepo := repository.NewTodoRepositoryMock()
			todoCache := cache.NewRedisCacheMock()

			todoRepo.On("Update", testCase.input).Return(testCase.repoUpdateReturn)
			if testCase.repoUpdateReturn == nil {

				todoCache.On("Get", mock.Anything, "todos").Return(testCase.cacheGetReturn.data, testCase.cacheGetReturn.err)

				if testCase.cacheGetReturn.err != nil {
					todoRepo.On("FindAll").Return(testCase.repoFindAllReturn.todos, testCase.repoFindAllReturn.err)
				}
				todoCache.On("Set", mock.Anything, "todos", mock.Anything, mock.Anything).Return(testCase.cacheSetReturn)
			}

			todoService := service.NewTodoService(todoRepo, todoCache)

			// Act
			err := todoService.Update(testCase.input)

			// Assert
			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			todoRepo.AssertExpectations(t)
		})
	}
}

func TestTodoserviceDelete(t *testing.T) {
	testCases := []struct {
		description       string
		input             dto.TodoInputDelete
		repoDeleteReturn  error
		repoFindAllReturn struct {
			todos []dto.Todo
			err   error
		}
		cacheGetReturn struct {
			data string
			err  error
		}
		cacheSetReturn error
		expectedErr    error
	}{
		{
			description: "Delete success",
			input: dto.TodoInputDelete{
				Id: "1",
			},
			repoDeleteReturn: nil,
			repoFindAllReturn: struct {
				todos []dto.Todo
				err   error
			}{[]dto.Todo{}, nil},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "[{\"id\":\"1\",\"topic\":\"Complete Project\",\"description\":\"Description for Complete Project\",\"status\":\"Completed\"}]",
				err:  nil,
			},
			cacheSetReturn: nil,
			expectedErr:    nil,
		},
		{
			description: "failed repository Delete.",
			input: dto.TodoInputDelete{
				Id: "1",
			},
			repoDeleteReturn: errors.New("failed repository delete."),
			repoFindAllReturn: struct {
				todos []dto.Todo
				err   error
			}{[]dto.Todo{}, nil},
			cacheGetReturn: struct {
				data string
				err  error
			}{data: "", err: nil},
			cacheSetReturn: nil,
			expectedErr:    errors.New("failed repository delete."),
		},
		{
			description: "Cache hit",
			input: dto.TodoInputDelete{
				Id: "1",
			},
			repoDeleteReturn: nil,
			repoFindAllReturn: struct {
				todos []dto.Todo
				err   error
			}{[]dto.Todo{}, nil},
			cacheGetReturn: struct {
				data string
				err  error
			}{
				data: "[{\"id\":\"1\",\"topic\":\"Complete Project\",\"description\":\"Description for Complete Project\",\"status\":\"Completed\"}]",
				err:  nil,
			},
			cacheSetReturn: nil,
			expectedErr:    nil,
		},
		{
			description: "Cache miss.",
			input: dto.TodoInputDelete{
				Id: "1",
			},
			repoDeleteReturn: nil,
			repoFindAllReturn: struct {
				todos []dto.Todo
				err   error
			}{
				[]dto.Todo{
					{
						Id:          "1",
						Topic:       "Complete Project",
						Description: "Description for Complete Project",
						Status:      "Completed",
					},
				}, nil},
			cacheGetReturn: struct {
				data string
				err  error
			}{data: "", err: errors.New("cache miss")},
			cacheSetReturn: nil,
			expectedErr:    nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			// Arrange
			todoRepo := repository.NewTodoRepositoryMock()
			todoCache := cache.NewRedisCacheMock()

			todoRepo.On("Delete", testCase.input).Return(testCase.repoDeleteReturn)

			if testCase.repoDeleteReturn == nil {
				todoCache.On("Get", mock.Anything, "todos").Return(testCase.cacheGetReturn.data, testCase.cacheGetReturn.err)
				if testCase.cacheGetReturn.err != nil {
					todoRepo.On("FindAll").Return(testCase.repoFindAllReturn.todos, testCase.repoFindAllReturn.err)
				}
				todoCache.On("Set", mock.Anything, "todos", mock.Anything, mock.Anything).Return(testCase.cacheSetReturn)
			}

			todoService := service.NewTodoService(todoRepo, todoCache)

			// Act
			err := todoService.Delete(testCase.input)

			// Assert
			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			todoRepo.AssertExpectations(t)
		})
	}
}
