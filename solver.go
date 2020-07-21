package main

import (
	"math/bits"
	"errors"
	"fmt"
	"log"
	"strconv"
	"unicode"
)

const L = 9 // Length of each side of the board
const B = 3 // Length of each side of a block
const BOARD_SIZE = L * L // Number of slots on the board
const BLOCK_SIZE = B * B // Number of slots on a block
const ONES = 0x1ff // 9-bit ones

type Board [BOARD_SIZE]uint

type State struct {
	// The current board
	board *Board
	// Available numbers for each slot. Each bit represents whether
	// the corresponding number is available due to the number not
	// taken by row, col, or block. Lowest bit is for 1, and
	// highest bit is for 9.
	avail *Board
	// Required numbers for each slot. Each bit represent whether
	// the corresponding number is the required number due to
	// all other slots in the block with the number not available.
	// Lowest bit is for 1, and highest bit is for 9.
	must *Board
}

func NewState() *State {
	retv := new(State)
	retv.board = new(Board)
	retv.avail = new(Board)
	retv.must = new(Board)
	
	for i, _ := range retv.avail {
		retv.avail[i] = ONES
		retv.must[i] = ONES
	}
	return retv
}

func (me *State) Set(row, col int, num uint) {
	pos := toPos(row, col)
	me.board[pos] = num
	me.updateAvail(row, col)
	me.updateMust(row, col)
}

func (me *State) updateAvail(row, col int) {
	pos := toPos(row, col)
	mask := ^(uint(1) << (me.board[pos] - 1))

	// update the slot
	me.avail[pos] = ^mask

	for i := 0; i < L; i++ {
		// update for row
		if col != i {
			me.avail[toPos(row, i)] &= mask
		}
		// update for col
		if row != i {
			me.avail[toPos(i, col)] &= mask
		}
	}
	// update for 3x3 block
	startRow := toBlockStart(row)
	startCol := toBlockStart(col)

	for r := 0; r < B; r++ {
		for c := 0; c < B; c++ {
			if startRow+r != row && startCol+c != col {
				me.avail[toPos(startRow+r, startCol+c)] &= mask
			}
		}
	}
}

func (me *State) updateMust(row, col int) {
	startRow := toBlockStart(row)
	startCol := toBlockStart(col)
	for r := 0; r < L; r += 3 {
		for c := 0; c < L; c += 3 {
			if (r == startRow || c == startCol) {
				me.updateMustBlock(r, c)
			}
		}
	}
}

func (me *State) updateMustBlock(blockStartRow, blockStartCol int) {
	var availBefore [BLOCK_SIZE + 2]uint // The logical OR of all avail value before (one-indexed)
	var availAfter [BLOCK_SIZE + 2]uint // The logical OR of all avail value after (one-indexed)
	for r := 0; r < B; r++ {
		for c := 0; c < B; c++ {
			blockPos := r * B + c + 1
			availBefore[blockPos] = availBefore[blockPos-1] | me.avail[toPos(blockStartRow+r, blockStartCol+c)]
		}
	}
	for r := B-1; r >= 0; r-- {
		for c := B-1; c >= 0; c-- {
			blockPos := r * B + c + 1
			availAfter[blockPos] = availAfter[blockPos+1] | me.avail[toPos(blockStartRow+r, blockStartCol+c)]
		}
	}
	for r := 0; r < B; r++ {
		for c := 0; c < B; c++ {
			blockPos := r * B + c + 1
			me.must[toPos(blockStartRow+r, blockStartCol+c)] = ONES - (availBefore[blockPos-1] | availAfter[blockPos+1])
		}
	}
}

func (me *State) Board() *Board {
	return me.board
}

func toPos(row, col int) int {
	return row * L + col
}

func toRowCol(pos int) (int, int) {
	return pos / L, pos % L
}

func toBlockStart(rowOrCol int) int {
	return (rowOrCol / B) * B
}

func readBoard() (*State, error) {
	retv := NewState()

	for row := 0; row < L; row++ {
		var line string
		_, err := fmt.Scanln(&line)
		if err != nil {
			return nil, err
		}
		// to take account of the \n
		if len(line) != L {
			return nil, errors.New("Each input line must have 9 characters.")
		}
		for col, c := range line {
			if unicode.IsDigit(c) {
				i, _ := strconv.Atoi(string(c))
				retv.Set(row, col, uint(i))
			} else if c != '.' {
				return nil, errors.New("Only digits and . is allowed for input")
			}
		}
	}
	return retv, nil
}

func writeBoard(board *Board) {
	for i, num := range board {
		if i != 0 && i%L == 0 {
			fmt.Println()
		}
		if num == 0 {
			fmt.Print(".")
		} else {
			fmt.Print(num)
		}
	}
	fmt.Println()
}

func solve(state *State) {
	for i := 0; i < BOARD_SIZE; i++ {
		row, col := toRowCol(i)
		if bits.OnesCount(state.avail[i]) == 1 {
			state.Set(row, col, uint(bits.TrailingZeros(state.avail[i]) + 1))
		} else if bits.OnesCount(state.must[i]) == 1 {
			state.Set(row, col, uint(bits.TrailingZeros(state.must[i]) + 1))
		}
	}
}

func main() {
	state, err := readBoard()
	if err != nil {
		log.Fatal(err)
	}
	solve(state)
	writeBoard(state.Board())
}