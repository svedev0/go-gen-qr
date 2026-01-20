package main

import "fmt"

// getNumRawDataModules returns the number of modules available for codewords
// (data + ECC), excluding all function patterns. This formula is a well-known
// closed form.
func getNumRawDataModules(version int) int {
	if version < 1 || version > 40 {
		panic("version out of range")
	}

	var result int = (16*version+128)*version + 64
	if version >= 2 {
		var numAlign int = version/7 + 2
		result -= (25*numAlign-10)*numAlign - 55
		if version >= 7 {
			result -= 36 // version information areas
		}
	}

	return result
}

func getNumDataCodewords(version int) int {
	var rawCodewords int = getNumRawDataModules(version) / 8
	var ecc int = eccCodewordsPerBlockL[version] * numEccBlocksL[version]
	return rawCodewords - ecc
}

func charCountBitsForByteMode(version int) int {
	if version <= 9 {
		return 8
	}
	return 16 // for 10..40
}

func chooseVersionAndData(payload []byte) (version int, dataCodewords []byte, err error) {
	if len(payload) == 0 {
		payload = []byte{}
	}

	for ver := 1; ver <= 40; ver++ {
		var dataCwCount int = getNumDataCodewords(ver)
		codewords, ok := buildDataCodewordsByte(payload, ver, dataCwCount)
		if ok {
			return ver, codewords, nil
		}
	}

	return 0, nil, fmt.Errorf("content too large: %d bytes", len(payload))
}

func buildDataCodewordsByte(payload []byte, version int, dataCwCount int) (codewords []byte, ok bool) {
	var bb bitBuf

	bb.put(modeByte, 4)
	bb.put(len(payload), charCountBitsForByteMode(version))
	for _, b := range payload {
		bb.put(int(b), 8)
	}

	var maxBits int = dataCwCount * 8
	if bb.len() > maxBits {
		return nil, false
	}

	// Terminator (up to 4 zeros)
	var rem int = maxBits - bb.len()
	var term int = min(rem, 4)
	bb.put(0, term)

	// Pad to byte boundary
	for bb.len()%8 != 0 {
		bb.put(0, 1)
	}

	// Convert to codewords
	var cw []byte = bb.bytes()
	if len(cw) > dataCwCount {
		return nil, false
	}

	// Pad codewords
	var pads []byte = []byte{0xEC, 0x11}
	for i := 0; len(cw) < dataCwCount; i++ {
		cw = append(cw, pads[i%2])
	}

	return cw, true
}

func codewordsToBits(codewords []byte) []bool {
	var bits []bool = make([]bool, 0, len(codewords)*8)
	for _, b := range codewords {
		for i := 7; i >= 0; i-- {
			bits = append(bits, ((b>>i)&1) == 1)
		}
	}
	return bits
}
