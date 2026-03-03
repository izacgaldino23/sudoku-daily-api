package vo

import "math/bits"

type Binary uint16

func NewBinary(n int) Binary {
	return Binary(1 << n)
}

func NewFullBinary(max int) Binary {
	return Binary((1 << (max + 1)) - 2)
}

func (b *Binary) Add(other int) {
	*b |= (1 << other)
}

func (b *Binary) Union(other ...Binary) {
	for _, o := range other {
		*b |= o
	}
}

func (b *Binary) Remove(other int) {
	*b &^= (1 << other)
}

func (b *Binary) Contains(other int) bool {
	return (*b & (1 << other)) != 0
}

func (b *Binary) Missing(other Binary) Binary {
	return other &^ *b
}

func (b *Binary) Count() int {
	return bits.OnesCount16(uint16(*b))
}

func (b *Binary) Sub(other Binary) {
	*b &^= other
}

func (b *Binary) Values() []int {
	var result []int
	v := uint16(*b)

	for v != 0 {
		pos := bits.TrailingZeros16(v)
		result = append(result, pos)
		v &= v - 1 // remove o bit menos significativo ligado
	}

	return result
}