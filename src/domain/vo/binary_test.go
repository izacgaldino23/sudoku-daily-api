package vo

import (
	"reflect"
	"testing"
)

func TestBinary_Creation(t *testing.T) {
	t.Run("NewBinary deve ligar o n-ésimo bit", func(t *testing.T) {

		got := NewBinary(3)
		var want Binary = 8
		if got != want {
			t.Errorf("NewBinary(3) = %b; want %b", got, want)
		}
	})

	t.Run("NewFullBinary deve ligar bits de 1 até max", func(t *testing.T) {

		got := NewFullBinary(3)
		var want Binary = 14
		if got != want {
			t.Errorf("NewFullBinary(3) = %b; want %b", got, want)
		}
	})
}

func TestBinary_Operations(t *testing.T) {
	t.Run("Add e Contains", func(t *testing.T) {
		b := NewBinary(0)
		b.Add(4)

		if !b.Contains(0) || !b.Contains(4) {
			t.Errorf("Esperava conter bits 0 e 4. Estado: %b", b)
		}
		if b.Contains(1) {
			t.Error("Não esperava conter o bit 1")
		}
	})

	t.Run("Union", func(t *testing.T) {
		b := NewBinary(1)

		b.Union(NewBinary(2), NewBinary(5))

		if !b.Contains(1) || !b.Contains(2) || !b.Contains(5) {
			t.Errorf("Union falhou. Estado: %b", b)
		}
	})

	t.Run("Remove", func(t *testing.T) {
		b := NewFullBinary(3)
		b.Remove(2)

		if b.Contains(2) {
			t.Error("Bit 2 deveria ter sido removido")
		}
		if b.Count() != 2 {
			t.Errorf("Count esperado 2, obtido %d", b.Count())
		}
	})
}

func TestBinary_Missing(t *testing.T) {
	owned := NewBinary(1)
	required := NewFullBinary(4)

	got := owned.Missing(required)

	want := NewFullBinary(4)
	want.Remove(1)

	if got != want {
		t.Errorf("Missing = %b; want %b", got, want)
	}
}

func TestBinary_Values(t *testing.T) {
	b := NewFullBinary(3)
	b.Add(5)

	got := b.Values()
	want := []int{1, 2, 3, 5}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Values() = %v; want %v", got, want)
	}
}
