package main

import (
	"fmt"
	"github.com/tygermarshall/blunderbuss/shared/board"
	"github.com/tygermarshall/blunderbuss/shared/pieces"
	"strings"
)

func printBoard(b *board.Board) {
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			p := b.Squares[rank][file]
			if p.Type == pieces.Empty {
				fmt.Print(". ")
				continue
			}

			c := p.Type.String()[0]
			if p.Team == pieces.Black {
				c = byte(strings.ToLower(string(c))[0])
			}
			fmt.Printf("%c ", c)
		}
		fmt.Println()
	}
}

const (
	reset = "\033[0m"

	bgLight = "\033[48;5;250m" // light square
	bgDark  = "\033[48;5;240m" // dark square

	fgWhite = "\033[38;5;231m"
	fgBlack = "\033[38;5;16m"
)

func PrettyPrint(b board.Board) {
	for rank := 0; rank <= 7; rank++ {
		fmt.Printf("%d ", rank+1)

		for file := 0; file < 8; file++ {
			p := b.Squares[rank][file]

			bg := squareColor(rank, file)
			fmt.Print(bg)

			if p.Type == pieces.Empty {
				fmt.Print("  ")
			} else {
				fg := fgWhite
				if p.Team == pieces.Black {
					fg = fgBlack
				}
				fmt.Print(fg, pieceRune(p), " ")
			}

			fmt.Print(reset)
		}
		fmt.Println()
	}

	fmt.Println("  a b c d e f g h")
}
func squareColor(rank, file int) string {
	if (rank+file)%2 == 0 {
		return bgLight
	}
	return bgDark
}

func pieceRune(p pieces.Piece) string {
	switch p.Team {
	case pieces.White:
		switch p.Type {
		case pieces.Pawn:
			return "♙"
		case pieces.Knight:
			return "♘"
		case pieces.Bishop:
			return "♗"
		case pieces.Rook:
			return "♖"
		case pieces.Queen:
			return "♕"
		case pieces.King:
			return "♔"
		}
	case pieces.Black:
		switch p.Type {
		case pieces.Pawn:
			return "♟"
		case pieces.Knight:
			return "♞"
		case pieces.Bishop:
			return "♝"
		case pieces.Rook:
			return "♜"
		case pieces.Queen:
			return "♛"
		case pieces.King:
			return "♚"
		}
	}

	return "?"
}
