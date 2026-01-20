// Penalty scoring (mask evaluation)

package main

func penaltyScore(mx *matrix) int {
	return penaltyRule1(mx) +
		penaltyRule2(mx) +
		penaltyRule3(mx) +
		penaltyRule4(mx)
}

func penaltyRule1(mx *matrix) int {
	var score int = 0
	var size int = mx.size

	for row := range size {
		var runColor bool = mx.m[row][0].val
		var runLen int = 1
		for col := 1; col < size; col++ {
			if mx.m[row][col].val == runColor {
				runLen++
			} else {
				if runLen >= 5 {
					score += 3 + (runLen - 5)
				}
				runColor = mx.m[row][col].val
				runLen = 1
			}
		}
		if runLen >= 5 {
			score += 3 + (runLen - 5)
		}
	}

	for col := range size {
		var runColor bool = mx.m[0][col].val
		var runLen int = 1
		for row := 1; row < size; row++ {
			if mx.m[row][col].val == runColor {
				runLen++
			} else {
				if runLen >= 5 {
					score += 3 + (runLen - 5)
				}
				runColor = mx.m[row][col].val
				runLen = 1
			}
		}
		if runLen >= 5 {
			score += 3 + (runLen - 5)
		}
	}

	return score
}

func penaltyRule2(mx *matrix) int {
	var score int = 0
	var size int = mx.size
	for row := 0; row < size-1; row++ {
		for col := 0; col < size-1; col++ {
			var val bool = mx.m[row][col].val
			if mx.m[row][col+1].val == val &&
				mx.m[row+1][col].val == val &&
				mx.m[row+1][col+1].val == val {
				score += 3
			}
		}
	}
	return score
}

func penaltyRule3(mx *matrix) int {
	var score int = 0
	var size int = mx.size
	var pattern1 []bool = []bool{
		true, false, true, true, true, false, true, false, false, false, false}
	var pattern2 []bool = []bool{
		false, false, false, false, true, false, true, true, true, false, true}

	for row := range size {
		for col := 0; col <= size-11; col++ {
			if matchRow(mx, row, col, pattern1) ||
				matchRow(mx, row, col, pattern2) {
				score += 40
			}
		}
	}

	for col := range size {
		for row := 0; row <= size-11; row++ {
			if matchCol(mx, row, col, pattern1) ||
				matchCol(mx, row, col, pattern2) {
				score += 40
			}
		}
	}

	return score
}

func penaltyRule4(mx *matrix) int {
	var size int = mx.size
	var darkCount int = 0
	var totalCount int = size * size
	for row := range size {
		for col := range size {
			if mx.m[row][col].val {
				darkCount++
			}
		}
	}

	var darkPercent int = (darkCount * 100) / totalCount
	var lower5Percent int = (darkPercent / 5) * 5
	var upper5Percent int = lower5Percent + 5
	var lowerDeviation int = abs(lower5Percent - 50)
	var upperDeviation int = abs(upper5Percent - 50)
	var nearestDeviation int = min(upperDeviation, lowerDeviation)
	return (nearestDeviation / 5) * 10
}

func matchRow(mx *matrix, row, col int, pattern []bool) bool {
	for i := range pattern {
		if mx.m[row][col+i].val != pattern[i] {
			return false
		}
	}
	return true
}

func matchCol(mx *matrix, row, col int, pattern []bool) bool {
	for i := range pattern {
		if mx.m[row+i][col].val != pattern[i] {
			return false
		}
	}
	return true
}
