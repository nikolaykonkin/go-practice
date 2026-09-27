package mymath

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"положительные", 2, 3, 5},
		{"с нулём", 0, 5, 5},
		{"отрицательные", -3, -2, -5},
		{"смешанные", -3, 5, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.want {
				t.Errorf("Add(%d, %d) = %d, ожидалось %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"положительные", 3, 4, 12},
		{"с нулём", 5, 0, 0},
		{"отрицательные", -3, -2, 6},
		{"смешанные", -3, 5, -15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Multiply(tt.a, tt.b); got != tt.want {
				t.Errorf("Multiply(%d, %d) = %d, ожидалось %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
