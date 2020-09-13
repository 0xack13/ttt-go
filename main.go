package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
)

// Constants
const (
	AGENT = +1
	HUMAN = -1
	NONE  = 0
	BLANK = ' '
	LINE  = "\n---------------"
)

type pair struct {
	x int
	y int
}

var board [][]int
var blanks = [][]int{}

// eval
func eval(state [][]int) int {
	if winner(state, AGENT) {
		return +1
	} else if winner(state, HUMAN) {
		return -1
	} else {
		return 0
	}
}

// Print board
func printboard(board [][]int) {

	pieces := map[int]rune{
		1:  'X',
		-1: 'O',
		0:  BLANK,
	}
	fmt.Println(LINE)
	for _, row := range board {
		for _, col := range row {
			fmt.Printf("| %c |", pieces[col])
		}
		fmt.Println(LINE)
	}
}

// blank tiles
func blankTiles(state [][]int) [][]int {
	blanks = [][]int{}
	for i, row := range state {
		for j, tile := range row {
			if tile == 0 {
				blanks = append(blanks, []int{i, j})
			}
		}
	}
	return blanks
}

// game over
func gameOver(state [][]int) bool {
	return winner(state, AGENT) || winner(state, HUMAN)
}

// Winner board
func winner(state [][]int, player int) bool {
	winBoard := [][]int{
		{state[0][0], state[0][1], state[0][2]},
		{state[1][0], state[1][1], state[1][2]},
		{state[2][0], state[2][1], state[2][2]},
		{state[0][0], state[1][0], state[2][0]},
		{state[0][1], state[1][1], state[2][1]},
		{state[0][2], state[1][2], state[2][2]},
		{state[0][0], state[1][1], state[2][2]},
		{state[2][0], state[1][1], state[0][2]},
	}
	playerrow := []int{player, player, player}
	for _, row := range winBoard {
		if reflect.DeepEqual(playerrow, row) {
			return true
		}
	}
	return false
}

func validAction(i int, j int) bool {
	playerrow := []int{i, j}
	for _, row := range blankTiles(board) {
		if reflect.DeepEqual(row, playerrow) {
			return true
		}
	}
	return false
}

func applyAction(i int, j int, player int) bool {
	if validAction(i, j) {
		board[i][j] = player
		return true
	} else {
		return false
	}
}

func miniMax(depth int, player int) []int {
	const MaxUint = ^uint(0)
	const MinUint = 0
	const MaxInt = int(MaxUint >> 1)
	const MinInt = -MaxInt - 1
	var maximise []int
	if player == AGENT {
		maximise = []int{-1, -1, MinInt}
	} else {
		maximise = []int{-1, -1, MaxInt}
	}

	if depth == 0 || gameOver(board) {
		score := eval(board)
		return []int{-1, -1, score}
	}

	// playerrow := []int{i, j}
	for _, row := range blankTiles(board) {
		i, j := row[0], row[1]
		board[i][j] = player
		player1 := player
		player1 *= -1
		score := miniMax(depth-1, player1)
		board[i][j] = 0
		score[0], score[1] = i, j
		if player == AGENT {
			if score[2] > maximise[2] {
				maximise = score
			}
		} else {
			if score[2] < maximise[2] {
				maximise = score
			}
		}
	}
	return maximise
}

func main() {
	var loc int
	stdin := bufio.NewReader(os.Stdin)

	// Actions map
	actions := map[int]pair{
		0: {0, 0}, 1: {0, 1}, 2: {0, 2},
		3: {1, 0}, 4: {1, 1}, 5: {1, 2},
		6: {2, 0}, 7: {2, 1}, 8: {2, 2},
	}

	board = [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	player := 1

	for len(blankTiles(board)) > 0 && !gameOver(board) {
		if player < 0 {
			fmt.Fscan(stdin, &loc)
			coord := actions[loc]
			if validAction(coord.x, coord.y) {
				player *= -1
				applyAction(coord.x, coord.y, -1)
			}
		} else {
			move := miniMax(len(blankTiles(board)), AGENT)
			// fmt.Println("%d %d", move[0], move[1])
			player *= -1
			applyAction(move[0], move[1], AGENT)
		}
		print("\033[H\033[2J")
		printboard(board)
	}
}
