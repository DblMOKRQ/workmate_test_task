package storage

import (
	"sync"
	"time"

	"github.com/DblMOKRQ/workmate_test_task/internal/models"
	"github.com/google/uuid"
)

type TaskStorage struct {
	sync.Mutex
	tasks map[string]*models.Task
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks: make(map[string]*models.Task),
	}
}

func (s *TaskStorage) GetTask(id string) (*models.Task, error) {
	s.Lock()
	defer s.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return nil, models.ErrTaskNotFound
	}
	return task, nil
}

func (s *TaskStorage) CreateTask() *models.Task {
	s.Lock()
	defer s.Unlock()

	id := generateID()
	task := &models.Task{
		ID:        id,
		Status:    models.TaskStatusPending,
		CreatedAt: time.Now(),
	}
	s.tasks[id] = task
	return task
}

func (s *TaskStorage) DeleteTask(id string) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return models.ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}

func generateID() string {
	return uuid.New().String()
}
