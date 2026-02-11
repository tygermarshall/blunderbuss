package main

import (
	"github.com/tygermarshall/blunderbuss/shared/board"
)

type Game struct {
	Board      board.Board
	TurnNumber int
}

func newGame() Game {
	return Game{Board: board.CreateDefaultBoard(),
		TurnNumber: 1,
	}
}
