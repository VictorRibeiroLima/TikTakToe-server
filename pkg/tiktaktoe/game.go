package tiktaktoe

import (
	"errors"
	"fmt"
)

type game struct {
	table [3][3]int8
	turn  int8
	moves int8
}

func (g game) Draw() {
	for i := 0; i < len(g.table); i++ {
		for j := 0; j < len(g.table); j++ {
			var symbol string
			square := g.table[i][j]
			if square == 0 {
				symbol = " "
			} else if square == 1 {
				symbol = "x"
			} else {
				symbol = "o"
			}
			fmt.Print(symbol)
			if j != 2 {
				fmt.Print(" | ")
			}
		}
		fmt.Println()
		if i != 2 {
			fmt.Println("----------")
		}
	}
}

func (g game) checkWinner(row int8, column int8) int8 {
	length := len(g.table)
	//checking column
	for i := 0; i < length; i++ {
		if g.table[row][i] != g.turn {
			break
		}
		if i == length-1 {
			return 1
		}
	}

	//checking row
	for i := 0; i < length; i++ {
		if g.table[i][column] != g.turn {
			break
		}
		if i == length-1 {
			return 1
		}
	}

	//checking diag
	if row == column {
		//we're on a diagonal
		for i := 0; i < length; i++ {
			if g.table[i][i] != g.turn {
				break
			}
			if i == length-1 {
				return 1
			}
		}
	}

	//checking anti diag
	if row+column == int8(length)-1 {
		//we're on a diagonal
		for i := 0; i < length; i++ {
			if g.table[i][(length-1)-i] != g.turn {
				break
			}
			if i == length-1 {
				return 1
			}
		}
	}
	if g.moves == 9 {
		return 2
	}
	return 0
}

func (g *game) MakePlay(row int8, column int8) (int8, error) {
	var result int8
	if g.turn == 0 {
		g.turn = 1
	}
	if row < 0 || row > 2 {
		return result, errors.New("Invalid row")
	}
	if column < 0 || column > 2 {
		return result, errors.New("Invalid column")
	}
	g.moves++
	if g.table[row][column] != 0 {
		return result, errors.New("square already marked")
	}
	g.table[row][column] = g.turn
	result = g.checkWinner(row, column)
	if g.turn == 1 {
		g.turn = 2
	} else {
		g.turn = 1
	}
	return result, nil
}
