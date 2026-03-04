package vo

import (
	"reflect"
	"testing"
)

func TestNewBinary(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected Binary
	}{
		{"sets bit 0", 0, Binary(1)},
		{"sets bit 1", 1, Binary(2)},
		{"sets bit 3", 3, Binary(8)},
		{"sets bit 9", 9, Binary(512)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBinary(tt.n)
			if got != tt.expected {
				t.Errorf("NewBinary(%d) = %b; want %b", tt.n, got, tt.expected)
			}
		})
	}
}

func TestNewFullBinary(t *testing.T) {
	tests := []struct {
		name     string
		max      int
		expected Binary
	}{
		{"max 1", 1, Binary(2)},    // 0b10
		{"max 3", 3, Binary(14)},   // 0b1110
		{"max 4", 4, Binary(30)},   // 0b11110
		{"max 9", 9, Binary(1022)}, // 0b1111111110
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFullBinary(tt.max)
			if got != tt.expected {
				t.Errorf("NewFullBinary(%d) = %b; want %b", tt.max, got, tt.expected)
			}
		})
	}
}

func TestBinary_Add(t *testing.T) {
	b := Binary(0)
	b.Add(1)
	b.Add(3)
	b.Add(5)

	if !b.Contains(1) || !b.Contains(3) || !b.Contains(5) {
		t.Errorf("Expected bits 1, 3, 5 to be set. Got: %b", b)
	}
	if b.Contains(2) || b.Contains(4) {
		t.Errorf("Expected bits 2, 4 to not be set. Got: %b", b)
	}
}

func TestBinary_Union(t *testing.T) {
	b := NewBinary(1)
	b.Union(NewBinary(2), NewBinary(5))

	if !b.Contains(1) || !b.Contains(2) || !b.Contains(5) {
		t.Errorf("Expected bits 1, 2, 5 to be set. Got: %b", b)
	}
}

func TestBinary_Remove(t *testing.T) {
	b := NewFullBinary(5)
	b.Remove(1)
	b.Remove(3)

	if b.Contains(1) || b.Contains(3) {
		t.Errorf("Expected bits 1, 3 to be removed. Got: %b", b)
	}
	if !b.Contains(2) || !b.Contains(4) || !b.Contains(5) {
		t.Errorf("Expected bits 2, 4, 5 to remain. Got: %b", b)
	}
}

func TestBinary_Contains(t *testing.T) {
	b := NewFullBinary(4)

	tests := []struct {
		bit      int
		expected bool
	}{
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{0, false},
		{5, false},
		{9, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := b.Contains(tt.bit)
			if got != tt.expected {
				t.Errorf("Contains(%d) = %v; want %v", tt.bit, got, tt.expected)
			}
		})
	}
}

func TestBinary_Missing(t *testing.T) {
	owned := NewBinary(1)
	full := NewFullBinary(4)

	got := owned.Missing(full)

	want := NewFullBinary(4)
	want.Remove(1)
	if got != want {
		t.Errorf("Missing() = %b; want %b", got, want)
	}

	empty := Binary(0)
	missing := empty.Missing(full)
	if missing != full {
		t.Errorf("Missing() from empty = %b; want %b", missing, full)
	}
}

func TestBinary_Count(t *testing.T) {
	tests := []struct {
		name     string
		binary   Binary
		expected int
	}{
		{"empty", Binary(0), 0},
		{"single bit", NewBinary(3), 1},
		{"full 4", NewFullBinary(4), 4},
		{"full 9", NewFullBinary(9), 9},
		{"mixed", NewBinary(1) | NewBinary(3) | NewBinary(5), 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.binary.Count()
			if got != tt.expected {
				t.Errorf("Count() = %d; want %d", got, tt.expected)
			}
		})
	}
}

func TestBinary_Sub(t *testing.T) {
	b := NewFullBinary(4)
	other := NewBinary(1) | NewBinary(2)

	b.Sub(other)

	if b.Contains(1) || b.Contains(2) {
		t.Errorf("Expected bits 1, 2 to be removed. Got: %b", b)
	}
	if !b.Contains(3) || !b.Contains(4) {
		t.Errorf("Expected bits 3, 4 to remain. Got: %b", b)
	}
}

func TestBinary_Values(t *testing.T) {
	tests := []struct {
		name     string
		binary   Binary
		expected []int
	}{
		{"empty", Binary(0), nil},
		{"single", NewBinary(5), []int{5}},
		{"full 3", NewFullBinary(3), []int{1, 2, 3}},
		{"full 4", NewFullBinary(4), []int{1, 2, 3, 4}},
		{"sparse", NewBinary(1) | NewBinary(4) | NewBinary(9), []int{1, 4, 9}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.binary.Values()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Values() = %v; want %v", got, tt.expected)
			}
		})
	}
}
