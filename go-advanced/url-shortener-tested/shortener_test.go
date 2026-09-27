package main

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestURLShortener_Shorten(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantErrIs error
	}{
		{"валидный HTTP URL", "http://example.com", nil},
		{"валидный HTTPS URL", "https://google.com/search?q=test", nil},
		{"невалидный URL", "not-a-url", ErrInvalidURL},
		{"пустая строка", "", ErrInvalidURL},
		{"URL без схемы", "example.com", ErrInvalidURL},
		{"URL с неподдерживаемой схемой", "ftp://example.com/file", ErrInvalidURL},
	}

	shortener := NewURLShortener()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortID, err := shortener.Shorten(tt.url)

			if tt.wantErrIs == nil {
				if err != nil {
					t.Fatalf("не ожидалась ошибка, получена: %v", err)
				}
				if len(shortID) < 6 || len(shortID) > 8 {
					t.Errorf("длина short ID = %d, ожидалось 6-8 (%s)", len(shortID), shortID)
				}
				return
			}

			if !errors.Is(err, tt.wantErrIs) {
				t.Errorf("ошибка = %v, ожидалась %v", err, tt.wantErrIs)
			}
		})
	}
}

func TestURLShortener_Shorten_Uniqueness(t *testing.T) {
	shortener := NewURLShortener()

	ids := make(map[string]bool)
	for i := 0; i < 20; i++ {
		url := fmt.Sprintf("https://example.com/%d", i)
		id, err := shortener.Shorten(url)
		if err != nil {
			t.Fatalf("неожиданная ошибка для %s: %v", url, err)
		}
		if ids[id] {
			t.Errorf("дублирующийся short ID: %s", id)
		}
		ids[id] = true
	}
}

func TestURLShortener_GetOriginal(t *testing.T) {
	shortener := NewURLShortener()

	originalURL := "https://example.com/very/long/path"
	shortID, err := shortener.Shorten(originalURL)
	if err != nil {
		t.Fatalf("не удалось подготовить данные: %v", err)
	}

	tests := []struct {
		name      string
		shortID   string
		want      string
		wantErrIs error
	}{
		{"существующий short_url", shortID, originalURL, nil},
		{"несуществующий short_url", "notexist", "", ErrShortNotFound},
		{"пустой short_url", "", "", ErrShortNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shortener.GetOriginal(tt.shortID)

			if tt.wantErrIs != nil {
				if !errors.Is(err, tt.wantErrIs) {
					t.Errorf("ошибка = %v, ожидалась %v", err, tt.wantErrIs)
				}
				return
			}

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if got != tt.want {
				t.Errorf("получено %q, ожидалось %q", got, tt.want)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"валидный http", "http://example.com", true},
		{"валидный https", "https://example.com/path?query=1", true},
		{"пустая строка", "", false},
		{"без схемы", "example.com", false},
		{"неподдерживаемая схема", "ftp://example.com", false},
		{"мусорная строка", "not-a-url", false},
		{"только схема без хоста", "http://", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidURL(tt.url)
			if got != tt.want {
				t.Errorf("isValidURL(%q) = %v, ожидалось %v", tt.url, got, tt.want)
			}
		})
	}
}

func TestGenerateShortID(t *testing.T) {
	for i := 0; i < 10; i++ {
		id, err := generateShortID()
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}

		if len(id) < 6 || len(id) > 8 {
			t.Errorf("длина = %d, ожидалось 6-8 (%s)", len(id), id)
		}

		for _, ch := range id {
			isAllowed := (ch >= 'A' && ch <= 'Z') ||
				(ch >= 'a' && ch <= 'z') ||
				(ch >= '0' && ch <= '9') ||
				ch == '-' || ch == '_'
			if !isAllowed {
				t.Errorf("недопустимый символ %q в ID %q", ch, id)
			}
		}
	}
}

func TestURLShortener_ConcurrentAccess(t *testing.T) {
	us := NewURLShortener()

	const workers = 10
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(workers * 2)

	for w := 0; w < workers; w++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_, _ = us.Shorten(fmt.Sprintf("https://example.com/%d/%d", id, i))
			}
		}(w)

		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_, _ = us.GetOriginal("nonexistent")
			}
		}()
	}

	wg.Wait()
}
