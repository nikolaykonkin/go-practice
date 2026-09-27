package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantURL    string
	}{
		{
			name:       "успешное сокращение",
			method:     http.MethodPost,
			body:       `{"url": "https://example.com/very/long/path"}`,
			wantStatus: http.StatusOK,
			wantURL:    "https://example.com/very/long/path",
		},
		{
			name:       "некорректный JSON",
			method:     http.MethodPost,
			body:       `{invalid json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "невалидный URL",
			method:     http.MethodPost,
			body:       `{"url": "not-a-url"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "пустой URL",
			method:     http.MethodPost,
			body:       `{"url": ""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "неподдерживаемый метод",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := NewURLShortener()
			handler := shortenHandler(us)

			req := httptest.NewRequest(tt.method, "/shorten", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("статус = %d, ожидался %d", rec.Code, tt.wantStatus)
			}

			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, ожидался application/json", ct)
			}

			if tt.wantStatus == http.StatusOK {
				var resp shortenResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("не удалось распарсить ответ: %v", err)
				}
				if len(resp.ShortURL) < 6 || len(resp.ShortURL) > 8 {
					t.Errorf("short_url должен быть 6-8 символов, получено: %s", resp.ShortURL)
				}
				if resp.OriginalURL != tt.wantURL {
					t.Errorf("original_url = %q, ожидался %q", resp.OriginalURL, tt.wantURL)
				}
				return
			}

			var errResp errorResponse
			if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
				t.Fatalf("не удалось распарсить тело ошибки: %v", err)
			}
			if errResp.Error == "" {
				t.Error("поле error в теле ответа пустое")
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	us := NewURLShortener()
	originalURL := "https://example.com/very/long/path"
	shortID, err := us.Shorten(originalURL)
	if err != nil {
		t.Fatalf("не удалось подготовить данные: %v", err)
	}

	tests := []struct {
		name         string
		method       string
		path         string
		wantStatus   int
		wantLocation string
	}{
		{
			name:         "успешный редирект",
			method:       http.MethodGet,
			path:         "/" + shortID,
			wantStatus:   http.StatusFound,
			wantLocation: originalURL,
		},
		{
			name:       "несуществующий short_url",
			method:     http.MethodGet,
			path:       "/doesnotexist",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "неподдерживаемый метод",
			method:     http.MethodPost,
			path:       "/" + shortID,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "пустой путь",
			method:     http.MethodGet,
			path:       "/",
			wantStatus: http.StatusNotFound,
		},
	}

	handler := redirectHandler(us)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("статус = %d, ожидался %d", rec.Code, tt.wantStatus)
			}

			if tt.wantLocation != "" {
				gotLocation := rec.Header().Get("Location")
				if gotLocation != tt.wantLocation {
					t.Errorf("Location = %q, ожидался %q", gotLocation, tt.wantLocation)
				}
				return
			}

			if tt.wantStatus >= 400 {
				var errResp errorResponse
				if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
					t.Fatalf("не удалось распарсить тело ошибки: %v", err)
				}
				if errResp.Error == "" {
					t.Error("поле error в теле ответа пустое")
				}
			}
		})
	}
}

// TestRouting проверяет регистрацию маршрутов через http.ServeMux,
// а не только сами обработчики напрямую
func TestRouting(t *testing.T) {
	us := NewURLShortener()
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", shortenHandler(us))
	mux.HandleFunc("/", redirectHandler(us))

	t.Run("POST /shorten через mux", func(t *testing.T) {
		body := bytes.NewBufferString(`{"url": "https://example.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/shorten", body)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("статус = %d, ожидался 200", rec.Code)
		}
	})

	t.Run("GET /{short_id} через mux", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/doesnotexist", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("статус = %d, ожидался 404", rec.Code)
		}
	})
}
