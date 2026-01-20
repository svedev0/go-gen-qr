package main

type bitBuf struct {
	bits []bool
}

func (bb *bitBuf) put(value int, n int) {
	for i := n - 1; i >= 0; i-- {
		bb.bits = append(bb.bits, ((value>>i)&1) == 1)
	}
}

func (bb *bitBuf) len() int {
	return len(bb.bits)
}

func (bb *bitBuf) bytes() []byte {
	var out []byte = make([]byte, 0, (len(bb.bits)+7)/8)
	for i := 0; i < len(bb.bits); i += 8 {
		var val byte
		for j := range 8 {
			val <<= 1
			if i+j < len(bb.bits) && bb.bits[i+j] {
				val |= 1
			}
		}
		out = append(out, val)
	}
	return out
}
