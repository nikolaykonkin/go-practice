package storage

import (
	"errors"

	"github.com/nikolaykonkin/go-practice/tasks-api-crud/internal/models"
)

// ErrTaskNotFound возвращается, когда задача с указанным ID отсутствует
var ErrTaskNotFound = errors.New("task not found")

// Storage описывает контракт хранилища задач
// Любая реализация (in-memory, БД и т.д.) должна удовлетворять этому интерфейсу
type Storage interface {
	List() []models.Task
	Create(models.Task) (models.Task, error)
	Get(id int) (models.Task, bool)
	Update(id int, task models.Task) (models.Task, error)
	Delete(id int) error
}
