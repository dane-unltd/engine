package bitmask

type BitMask []byte

func New(nBits int) BitMask {
	nBytes := nBits >> 3
	if nBits&7 > 0 {
		nBytes++
	}
	return make(BitMask, nBytes)
}

func (b BitMask) Set(i int) {
	byteIx := i >> 3
	bitIx := i & 7
	b[byteIx] |= 1 << uint(bitIx)
}

func (b BitMask) UnSet(i int) {
	byteIx := i >> 3
	bitIx := i & 7
	b[byteIx] &= ^(1 << uint(bitIx))
}

func (b BitMask) Check(i int) bool {
	byteIx := i >> 3
	bitIx := i & 7
	return b[byteIx]&(1<<uint(bitIx)) > 0
}
