package main

type module struct {
	val      bool
	set      bool
	reserved bool
}

type matrix struct {
	size int
	m    [][]module
}

func newMatrix(size int) *matrix {
	var m [][]module = make([][]module, size)
	for i := range m {
		m[i] = make([]module, size)
	}
	return &matrix{size: size, m: m}
}

func (mx *matrix) clone() *matrix {
	var cp *matrix = newMatrix(mx.size)
	for row := 0; row < mx.size; row++ {
		copy(cp.m[row], mx.m[row])
	}
	return cp
}

func (mx *matrix) inBounds(row, col int) bool {
	return row >= 0 && row < mx.size && col >= 0 && col < mx.size
}

func (mx *matrix) setMod(row, col int, val bool, reserved bool) {
	if !mx.inBounds(row, col) {
		return
	}
	mx.m[row][col] = module{val: val, set: true, reserved: reserved}
}

func (mx *matrix) isData(row, col int) bool {
	return mx.m[row][col].set && !mx.m[row][col].reserved
}

// Function patterns
func placeFunctionPatterns(mx *matrix, version int) {
	var size int = mx.size

	// Finder patterns
	placeFinder(mx, 0, 0)
	placeFinder(mx, 0, size-7)
	placeFinder(mx, size-7, 0)

	// Separators
	placeSeparator(mx, 0, 0)
	placeSeparator(mx, 0, size-7)
	placeSeparator(mx, size-7, 0)

	// Timing patterns
	for i := 8; i <= size-9; i++ {
		mx.setMod(6, i, i%2 == 0, true)
		mx.setMod(i, 6, i%2 == 0, true)
	}

	// Alignment patterns
	var pos []int = getAlignmentPositions(version)
	for _, row := range pos {
		for _, col := range pos {
			// Skip the three alignment positions that overlap finder corners
			if (row == 6 && col == 6) ||
				(row == 6 && col == size-7) ||
				(row == size-7 && col == 6) {
				continue
			}
			placeAlignment(mx, row, col)
		}
	}

	// Dark module
	mx.setMod(4*version+9, 8, true, true)

	// Reserve format info areas
	reserveFormatInfo(mx)

	// Reserve version info areas for ver >= 7
	if version >= 7 {
		reserveVersionInfo(mx)
	}
}

func placeFinder(mx *matrix, row0, col0 int) {
	for row := range 7 {
		for col := range 7 {
			var finalRow, finalCol int = row0 + row, col0 + col
			var outer bool = row == 0 || row == 6 || col == 0 || col == 6
			var inner bool = row >= 2 && row <= 4 && col >= 2 && col <= 4
			mx.setMod(finalRow, finalCol, outer || inner, true)
		}
	}
}

func placeSeparator(mx *matrix, row0, col0 int) {
	for row := -1; row <= 7; row++ {
		for col := -1; col <= 7; col++ {
			var finalRow, finalCol int = row0 + row, col0 + col
			if !mx.inBounds(finalRow, finalCol) {
				continue
			}
			// Skip cells that are inside the finder itself
			if finalRow >= row0 &&
				finalRow <= row0+6 &&
				finalCol >= col0 &&
				finalCol <= col0+6 {
				continue
			}
			mx.setMod(finalRow, finalCol, false, true)
		}
	}
}

func placeAlignment(mx *matrix, row0, col0 int) {
	// 5x5 centered at (row0, col0). Black outer ring, white inner ring, black
	// center.
	for dy := -2; dy <= 2; dy++ {
		for dx := -2; dx <= 2; dx++ {
			var dist int = max(abs(dx), abs(dy))
			var isBlack bool = (dist == 2) || (dist == 0)
			mx.setMod(row0+dy, col0+dx, isBlack, true)
		}
	}
}

func getAlignmentPositions(version int) []int {
	if version == 1 {
		return nil
	}

	var size int = 17 + 4*version
	var numAlign int = version/7 + 2
	var pos []int = make([]int, numAlign)
	pos[0] = 6
	pos[numAlign-1] = size - 7

	if numAlign == 2 {
		return pos
	}

	var step int = (size - 13) / (numAlign - 1)
	if step%2 == 1 {
		step++
	}

	for i := numAlign - 2; i >= 1; i-- {
		pos[i] = pos[i+1] - step
	}

	return pos
}

func reserveFormatInfo(mx *matrix) {
	var size int = mx.size
	// Top-left format area: row 8 col 0..8 and col 8 row 0..8
	for col := 0; col <= 8; col++ {
		if !mx.m[8][col].set {
			mx.setMod(8, col, false, true)
		} else {
			mx.m[8][col].reserved = true
		}
	}

	for row := 0; row <= 8; row++ {
		if !mx.m[row][8].set {
			mx.setMod(row, 8, false, true)
		} else {
			mx.m[row][8].reserved = true
		}
	}

	// Other copy: row 8 col size-8..size-1 and col 8 row size-7..size-1
	for col := size - 8; col < size; col++ {
		if !mx.m[8][col].set {
			mx.setMod(8, col, false, true)
		} else {
			mx.m[8][col].reserved = true
		}
	}

	for row := size - 7; row < size; row++ {
		if !mx.m[row][8].set {
			mx.setMod(row, 8, false, true)
		} else {
			mx.m[row][8].reserved = true
		}
	}
}

func reserveVersionInfo(mx *matrix) {
	var size int = mx.size
	// Two 3x6 areas: top-right and bottom-left (reserved now, written later).
	for row := range 6 {
		for col := size - 11; col <= size-9; col++ {
			if !mx.m[row][col].set {
				mx.setMod(row, col, false, true)
			} else {
				mx.m[row][col].reserved = true
			}
		}
	}

	for row := size - 11; row <= size-9; row++ {
		for col := range 6 {
			if !mx.m[row][col].set {
				mx.setMod(row, col, false, true)
			} else {
				mx.m[row][col].reserved = true
			}
		}
	}
}

// Data placement
func placeDataBits(mx *matrix, bits []bool) {
	var size int = mx.size
	var i int = 0
	var up bool = true

	for col := size - 1; col >= 0; col -= 2 {
		if col == 6 {
			col--
		}

		var rows []int
		if up {
			rows = makeRange(size-1, 0, -1)
		} else {
			rows = makeRange(0, size-1, 1)
		}

		for _, row := range rows {
			for deltaCol := range 2 {
				var finalCol int = col - deltaCol
				if finalCol < 0 {
					continue
				}
				if mx.m[row][finalCol].set {
					continue
				}

				var bit bool = false
				if i < len(bits) {
					bit = bits[i]
					i++
				}

				mx.setMod(row, finalCol, bit, false)
			}
		}

		up = !up
	}
}

func makeRange(start, end, step int) []int {
	var output []int
	if step > 0 {
		for i := start; i <= end; i += step {
			output = append(output, i)
		}
	} else {
		for i := start; i >= end; i += step {
			output = append(output, i)
		}
	}
	return output
}

// Masking
func applyMask(mx *matrix, mask int) {
	for row := 0; row < mx.size; row++ {
		for col := 0; col < mx.size; col++ {
			if !mx.isData(row, col) {
				continue
			}
			if maskFunc(mask, row, col) {
				mx.m[row][col].val = !mx.m[row][col].val
			}
		}
	}
}

func maskFunc(mask, row, col int) bool {
	switch mask {
	case 0:
		return (row+col)%2 == 0
	case 1:
		return row%2 == 0
	case 2:
		return col%3 == 0
	case 3:
		return (row+col)%3 == 0
	case 4:
		return ((row/2)+(col/3))%2 == 0
	case 5:
		return (row*col)%2+(row*col)%3 == 0
	case 6:
		return ((row*col)%2+(row*col)%3)%2 == 0
	case 7:
		return ((row+col)%2+(row*col)%3)%2 == 0
	default:
		return false
	}
}
