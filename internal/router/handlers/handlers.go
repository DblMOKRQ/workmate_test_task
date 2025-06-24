package handlers

import (
	"encoding/json"
	"errors"

	"net/http"
	"time"

	"github.com/DblMOKRQ/workmate_test_task/internal/models"
	"github.com/DblMOKRQ/workmate_test_task/internal/storage"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	storage *storage.TaskStorage
	log     *zap.Logger
}

func NewHandler(storage *storage.TaskStorage, log *zap.Logger) *Handler {
	return &Handler{storage: storage, log: log.Named("handler")}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	task := h.storage.CreateTask()
	h.log.Info("Task created", zap.String("ID", task.ID))
	resp := models.CreateTaskResponse{TaskID: task.ID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	task, err := h.storage.GetTask(id)
	if err != nil {
		if errors.Is(err, models.ErrTaskNotFound) {
			h.log.Error("Task not found", zap.Error(err))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		h.log.Error("Error getting task", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.log.Info("Task retrieved", zap.String("ID", task.ID))

	// Рассчитываем продолжительность
	var duration time.Duration
	if !task.StartedAt.IsZero() {
		if !task.FinishedAt.IsZero() {
			duration = task.FinishedAt.Sub(task.StartedAt)
		} else {
			duration = time.Since(task.StartedAt)
		}
	}
	resp := models.TaskResponse{
		ID:           task.ID,
		Status:       task.Status,
		CreatedAt:    task.CreatedAt.Format(time.RFC3339),
		DurationSecs: duration.Seconds(),
	}
	if task.Status == "completed" {
		resp.Result = task.Result
	} else if task.Status == "failed" {
		resp.Error = task.Error
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DelteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.storage.DeleteTask(id)
	if err != nil {
		if errors.Is(err, models.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			h.log.Error("Task not found", zap.Error(err))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Error("Task deletion failed", zap.Error(err))
		json.NewEncoder(w).Encode(err)
		return
	}
	h.log.Info("Task deleted", zap.String("ID", id))
	w.WriteHeader(http.StatusNoContent)
}
