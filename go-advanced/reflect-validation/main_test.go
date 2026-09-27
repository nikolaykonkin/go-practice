package main

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   User
		wantErr string // подстрока в ошибке, пустая — ошибки нет
	}{
		{
			name:    "все поля валидны",
			input:   User{Name: "Иван", Age: 35, Email: "test@example.com"},
			wantErr: "",
		},
		{
			name:    "имя короче минимума",
			input:   User{Name: "Ив", Age: 35, Email: "test@example.com"},
			wantErr: "длина 2 меньше минимальной 3",
		},
		{
			name:    "возраст меньше минимума",
			input:   User{Name: "Иван", Age: 17, Email: "test@example.com"},
			wantErr: "значение 17 меньше минимального 18",
		},
		{
			name:    "возраст больше максимума",
			input:   User{Name: "Иван", Age: 70, Email: "test@example.com"},
			wantErr: "значение 70 больше максимального 65",
		},
		{
			name:    "email не соответствует regexp",
			input:   User{Name: "Иван", Age: 35, Email: "invalid email"},
			wantErr: "не соответствует формату",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("ожидалось nil, получено: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("ожидалась ошибка со словами %q, получено nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ошибка = %q, ожидалась подстрока %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestValidate_Nil(t *testing.T) {
	if err := Validate(nil); err == nil {
		t.Error("ожидалась ошибка для nil, получено nil")
	}
}

func TestValidate_NotStruct(t *testing.T) {
	if err := Validate(42); err == nil {
		t.Error("ожидалась ошибка для не-структуры, получено nil")
	}
}

func TestValidate_Pointer(t *testing.T) {
	u := &User{Name: "Иван", Age: 35, Email: "test@example.com"}
	if err := Validate(u); err != nil {
		t.Errorf("ожидалось nil для указателя на валидную структуру, получено: %v", err)
	}
}

func TestValidate_KirillicLength(t *testing.T) {
	// "Ив" — 2 руны, но 4 байта, валидатор должен считать руны
	err := Validate(User{Name: "Ив", Age: 35, Email: "test@example.com"})
	if err == nil {
		t.Fatal("ожидалась ошибка для имени из 2 рун")
	}
	if !strings.Contains(err.Error(), "длина 2") {
		t.Errorf("ожидалась длина 2, получено: %v", err)
	}
}
