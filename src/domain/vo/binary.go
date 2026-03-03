package vo

type Binary uint8

func NewBinary(n int) Binary {
	return Binary(1 << n)
}

func (b *Binary) Add(other int) {
	*b |= (1 << other)
}

func (b *Binary) Remove(other int) {
	*b &^= (1 << other)
}

func (b *Binary) Contains(other int) bool {
	return (*b & (1<<other)) != 0
}