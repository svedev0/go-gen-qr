package main

func writeFormatInfo(mx *matrix, mask int) {
	var size int = mx.size
	var format uint16 = formatBits(ecLevelL, mask) // 15 bits, MSB-first

	var coords1 = [15][2]int{
		{8, 0}, {8, 1}, {8, 2}, {8, 3}, {8, 4}, {8, 5},
		{8, 7},
		{8, 8},
		{7, 8},
		{5, 8}, {4, 8}, {3, 8}, {2, 8}, {1, 8}, {0, 8},
	}

	var coords2 [15][2]int
	// 7 bits: col 8, rows size - 1 .. size - 7
	for i := range 7 {
		coords2[i] = [2]int{size - 1 - i, 8}
	}

	// 8 bits: row 8, cols size - 8 .. size - 1
	for i := range 8 {
		coords2[7+i] = [2]int{8, size - 8 + i}
	}

	for i := range 15 {
		var bit bool = ((format >> (14 - i)) & 1) == 1
		var row1, col1 int = coords1[i][0], coords1[i][1]
		var row2, col2 int = coords2[i][0], coords2[i][1]
		mx.setMod(row1, col1, bit, true)
		mx.setMod(row2, col2, bit, true)
	}
}

func formatBits(ecBits int, mask int) uint16 {
	var data int = (ecBits << 3) | (mask & 0x7) // 5 bits
	var val uint16 = uint16(data) << 10
	var rem uint16 = bchRemainder(val, 0x537, 15) // generator 0x537, 15-bit field
	return ((uint16(data)<<10 | rem) ^ 0x5412) & 0x7FFF
}

func writeVersionInfo(mx *matrix, version int) {
	if version < 7 {
		return
	}
	var size int = mx.size
	var vbits uint32 = versionBits(version) // 18 bits

	for i := range 18 {
		var bit bool = ((vbits >> i) & 1) == 1 // LSB-first placement
		var row int = i / 3
		var col int = i % 3

		mx.setMod(row, size-11+col, bit, true) // Top-right
		mx.setMod(size-11+col, row, bit, true) // Bottom-left
	}
}

func versionBits(version int) uint32 {
	// 6-bit version number + 12-bit BCH (generator 0x1F25)
	var ver uint32 = uint32(version) << 12
	var rem uint32 = bchRemainder32(ver, 0x1F25, 18)
	return (uint32(version) << 12) | rem
}

func bchRemainder(val uint16, poly uint16, bits int) uint16 {
	var polyDeg int = degree16(poly)
	for i := bits - 1; i >= polyDeg; i-- {
		if (val>>i)&1 == 1 {
			val ^= poly << (i - polyDeg)
		}
	}
	return val & ((1 << polyDeg) - 1)
}

func bchRemainder32(val uint32, poly uint32, bits int) uint32 {
	var polyDeg int = degree32(poly)
	for i := bits - 1; i >= polyDeg; i-- {
		if (val>>i)&1 == 1 {
			val ^= poly << (i - polyDeg)
		}
	}
	return val & ((1 << polyDeg) - 1)
}

func degree16(num uint16) int {
	var deg int = -1
	for num != 0 {
		num >>= 1
		deg++
	}
	return deg
}

func degree32(num uint32) int {
	var deg int = -1
	for num != 0 {
		num >>= 1
		deg++
	}
	return deg
}
