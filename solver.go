package main

import (
	"math/bits"
	"errors"
	"fmt"
	"log"
	"strconv"
	"unicode"
)

const L = 9
const B = 3
const BOARD_SIZE = 81

type Board [BOARD_SIZE]uint

type State struct {
	// The current board
	board *Board
	// Available numbers for each slot. Each bit represents whether
	// the corresponding number is available. Lowest bit is for 1, and
	// highest bit is for 9.
	avail *Board
}

func NewState() *State {
	retv := new(State)
	retv.board = new(Board)
	retv.avail = new(Board)
	
	for i, _ := range retv.avail {
		retv.avail[i] = 0x1ff
	}
	return retv
}

func (me *State) Set(pos int, num uint) {
	me.board[pos] = num
	me.Update(pos)
}

func (me *State) Update(pos int) {
	row := pos / L
	col := pos % L
	mask := ^(uint(1) << (me.board[pos] - 1))
	for i := 0; i < L; i++ {
		// update for row
		if col != i {
			me.avail[row*L+i] &= mask
		}
		// update for col
		if row != i {
			me.avail[i*L+col] &= mask
		}
	}
	// update for 3x3 block
	startRow := (row / B) * B
	startCol := (col / B) * B

	for i := 0; i < B; i++ {
		for j := 0; j < B; j++ {
			if startRow+i != row && startCol+j != col {
				me.avail[(startRow+i)*L+startCol+j] &= mask
			}
		}
	}
}

func (me *State) Board() *Board {
	return me.board
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
				retv.Set(row*L+col, uint(i))
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
		if bits.OnesCount(state.avail[i]) == 1 {
			state.Set(i, uint(bits.TrailingZeros(state.avail[i]) + 1))
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
