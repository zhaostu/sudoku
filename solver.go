package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

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
				if err := retv.Set(row, col, uint(i)); err != nil {
					return nil, fmt.Errorf("Invalid input board, cannot set (%d %d) to %d", row, col, i)
				}
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

func solve(state *State) (*State, error) {
	for !state.Solved() {
		row, col, numbers := state.PickEmptySlot()
		if len(numbers) == 1 {
			if err := state.Set(row, col, numbers[0]); err != nil {
				return state, err
			}
		} else {
			// Here, we found there can be multiple numbers put into the slot
			for _, num := range numbers {
				stateCopy := *state
				if err := stateCopy.Set(row, col, num); err != nil {
					continue
				}
				stateSolved, err := solve(&stateCopy)
				if err == nil {
					return stateSolved, nil
				}
			}
			// When we reach here, all of the numbers has been tried and we found
			// no solution
			return state, errors.New("We did not find any solution to the board")
		}
	}
	return state, nil
}

func main() {
	state, err := readBoard()
	if err != nil {
		log.Fatal(err)
	}
	stateSolved, err := solve(state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error solving sudoku: %s. Board was last seen as:\n", err)
		writeBoard(stateSolved.Board())
	} else {
		writeBoard(stateSolved.Board())
	}
}
