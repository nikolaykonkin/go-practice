package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nikolaykonkin/go-practice/tasks-api-crud/internal/models"
	"github.com/nikolaykonkin/go-practice/tasks-api-crud/internal/storage"
)

// Handler содержит зависимости HTTP-обработчиков
type Handler struct{ Store storage.Storage }

// New создает новый Handler с переданным хранилищем
func New(s storage.Storage) *Handler { return &Handler{Store: s} }

// errorResponse — единый формат ошибок для всех эндпоинтов
type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON пишет JSON-ответ с нужным статусом и заголовком Content-Type
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// writeError пишет JSON-ответ с ошибкой в едином формате {"error": "..."}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

// decodeJSONBody строго разбирает JSON из тела запроса: отклоняет
// неизвестные поля и данные после первого JSON-документа
func decodeJSONBody(r *http.Request, v interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		return err
	}

	if dec.More() {
		return fmt.Errorf("тело запроса содержит лишние данные")
	}

	return nil
}

// TasksCollection обрабатывает /tasks: GET — список задач, POST — создание
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.List()
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := decodeJSONBody(r, &task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	created, err := h.Store.Create(task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// TaskItem обрабатывает /tasks/{id}: GET, PUT, DELETE
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, id)
	case http.MethodPut:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// extractID достает числовой ID из пути вида /tasks/{id}
// Путь должен содержать ровно два сегмента: "tasks" и числовой ID
// Пути вида /tasks/1/anything отклоняются как некорректные
func extractID(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 || parts[0] != "tasks" || parts[1] == "" {
		return 0, fmt.Errorf("invalid path")
	}
	return strconv.Atoi(parts[1])
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request, id int) {
	task, ok := h.Store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request, id int) {
	var task models.Task
	if err := decodeJSONBody(r, &task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	updated, err := h.Store.Update(id, task)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.Store.Delete(id); err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
