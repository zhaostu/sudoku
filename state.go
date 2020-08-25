package main

import (
	"fmt"
	"math/bits"
)

const L = 9              // Length of each side of the board
const B = 3              // Length of each side of a block
const BOARD_SIZE = L * L // Number of slots on the board
const BLOCK_SIZE = B * B // Number of slots on a block
const ONES = 0x1ff       // 9-bit ones

// Board is an array representing the board
type Board [BOARD_SIZE]uint

type State struct {
	// The current board
	board Board
	// Available numbers for each slot. Each bit represents whether
	// the corresponding number is available due to the number not
	// taken by row, col, or block. Lowest bit is for 1, and
	// highest bit is for 9.
	avail Board
	// Required numbers for each slot. Each bit represent whether
	// the corresponding number is the required number due to
	// all other slots in the block with the number not available.
	// Lowest bit is for 1, and highest bit is for 9.
	must   Board
	blanks int
}

func NewState() *State {
	retv := new(State)
	retv.blanks = BOARD_SIZE

	for i, _ := range retv.avail {
		retv.avail[i] = ONES
		retv.must[i] = ONES
	}
	return retv
}

func (me *State) Set(row, col int, num uint) error {
	pos := toPos(row, col)

	if me.board[pos] != 0 {
		return fmt.Errorf("Cannot set (%d %d) to %d, already set to %d.",
			row, col, num, me.board[pos])
	}
	me.board[pos] = num
	if err := me.updateAvail(row, col); err != nil {
		return err
	}
	me.updateMust(row, col)

	me.blanks -= 1
	return nil
}

func (me *State) updateAvail(row, col int) error {
	pos := toPos(row, col)
	mask := ^(uint(1) << (me.board[pos] - 1))

	// update the slot
	if me.avail[pos] & ^mask == 0 {
		return fmt.Errorf("Cannot set (%d %d) to %d, availability is %b",
			row, col, me.board[pos], me.avail[pos])
	}
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
	return nil
}

func (me *State) updateMust(row, col int) {
	startRow := toBlockStart(row)
	startCol := toBlockStart(col)
	for r := 0; r < L; r += 3 {
		for c := 0; c < L; c += 3 {
			if r == startRow || c == startCol {
				me.updateMustBlock(r, c)
			}
		}
	}
}

func (me *State) updateMustBlock(blockStartRow, blockStartCol int) {
	var availBefore [BLOCK_SIZE + 2]uint // The logical OR of all avail value before (one-indexed)
	var availAfter [BLOCK_SIZE + 2]uint  // The logical OR of all avail value after (one-indexed)
	for r := 0; r < B; r++ {
		for c := 0; c < B; c++ {
			blockPos := r*B + c + 1
			availBefore[blockPos] = availBefore[blockPos-1] | me.avail[toPos(blockStartRow+r, blockStartCol+c)]
		}
	}
	for r := B - 1; r >= 0; r-- {
		for c := B - 1; c >= 0; c-- {
			blockPos := r*B + c + 1
			availAfter[blockPos] = availAfter[blockPos+1] | me.avail[toPos(blockStartRow+r, blockStartCol+c)]
		}
	}
	for r := 0; r < B; r++ {
		for c := 0; c < B; c++ {
			blockPos := r*B + c + 1
			me.must[toPos(blockStartRow+r, blockStartCol+c)] = ONES - (availBefore[blockPos-1] | availAfter[blockPos+1])
		}
	}
}

// Returns a best column, row, and a list of possible values for that slot
func (me *State) PickEmptySlot() (int, int, []uint) {
	bestPos := -1
	bestAvailCount := 9
	for pos := 0; pos < BOARD_SIZE; pos++ {
		// Already filled
		if me.board[pos] != 0 {
			continue
		}
		row, col := toRowCol(pos)

		if me.must[pos] > 0 {
			return row, col, []uint{uint(bits.TrailingZeros(me.must[pos]) + 1)}
		}

		availCount := bits.OnesCount(me.avail[pos])
		if availCount == 1 {
			return row, col, []uint{uint(bits.TrailingZeros(me.avail[pos]) + 1)}
		}
		if availCount < bestAvailCount {
			bestPos = pos
			bestAvailCount = availCount
		}
	}
	bestRow, bestCol := toRowCol(bestPos)
	return bestRow, bestCol, toNumbers(me.avail[bestPos])
}

func (me *State) Board() *Board {
	return &me.board
}

func (me *State) Solved() bool {
	return me.blanks == 0
}
