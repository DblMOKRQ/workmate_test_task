package models

import (
	"errors"
	"time"
)

const (
	TaskStatusPending    = "pending"
	TaskStatusProcessing = "processing"
	TaskStatusCompleted  = "completed"
	TaskStatusFailed     = "failed"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	ID         string    // Уникальный идентификатор задачи
	Status     string    // "pending", "processing", "completed", "failed"
	CreatedAt  time.Time // Когда создана
	StartedAt  time.Time // Когда начали выполнять
	FinishedAt time.Time // Когда завершили
	Result     string    // Результат выполнения
	Error      error     // Ошибка (если была)
}

type TaskResponse struct {
	ID           string  `json:"id"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
	StartedAt    string  `json:"started_at,omitempty"`
	FinishedAt   string  `json:"finished_at,omitempty"`
	DurationSecs float64 `json:"duration_seconds,omitempty"`
	Result       string  `json:"result,omitempty"`
	Error        error   `json:"error,omitempty"`
}

type CreateTaskResponse struct {
	TaskID string `json:"task_id"`
}
