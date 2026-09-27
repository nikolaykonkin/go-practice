package models

// Task представляет одну задачу в списке дел
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at,omitempty"`
}
