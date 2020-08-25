package main

func toPos(row, col int) int {
	return row*L + col
}

func toRowCol(pos int) (int, int) {
	return pos / L, pos % L
}

func toBlockStart(rowOrCol int) int {
	return (rowOrCol / B) * B
}

func toNumbers(bits uint) []uint {
	numbers := make([]uint, 0, 9)
	current := uint(1)
	for bits > 0 {
		if bits&1 == 1 {
			numbers = append(numbers, current)
		}
		bits = bits >> 1
		current += 1
	}
	return numbers
}
