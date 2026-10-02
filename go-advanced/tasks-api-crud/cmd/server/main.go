package main

import (
	"log"
	"net/http"

	"github.com/nikolaykonkin/go-practice/tasks-api-crud/internal/handlers"
	"github.com/nikolaykonkin/go-practice/tasks-api-crud/internal/middleware"
	"github.com/nikolaykonkin/go-practice/tasks-api-crud/internal/storage"
)

func main() {
	store := storage.NewMemoryStorage()

	h := handlers.New(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.TasksCollection) // GET, POST
	mux.HandleFunc("/tasks/", h.TaskItem)       // GET, PUT, DELETE
	mux.HandleFunc("/health", healthHandler)    // GET
	mux.HandleFunc("/", notFoundHandler)        // любой другой путь → JSON 404

	handler := middleware.Logging(mux)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

// notFoundHandler возвращает JSON-ошибку для любого пути, не совпавшего
// ни с одним зарегистрированным маршрутом
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"not found"}`))
}

// healthHandler — простой health-check эндпоинт
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"method not allowed"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
