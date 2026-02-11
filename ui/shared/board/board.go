package board

import (
	"errors"
	"fmt"
	"github.com/tygermarshall/blunderbuss/shared/pieces"
)

type Board struct {
	Squares [8][8]pieces.Piece
}

type Coordinate struct {
	X int
	Y int
}

func (b Board) MovePiece(start, end Coordinate) (Board, error) {
	piece, err := b.getPiece(start)
	if err != nil {
		return b, err
	}
	switch piece.Type {
	case pieces.Pawn:
		b, err = b.movePawn(start, end, &piece)
		if err != nil {
			return b, err
		}
		return b, nil

	case pieces.Rook:
		b, err = b.moveRook(start, end, piece)
		return b, nil
	case pieces.Knight:
		b, err = b.moveKnight(start, end, piece)
		return b, nil
	case pieces.Bishop:
		b, err = b.moveBishop(start, end, piece)
		return b, nil
	case pieces.Queen:
		b, err = b.moveQueen(start, end, piece)
		return b, nil
	case pieces.King:
		b, err = b.moveKing(start, end, piece)
		return b, nil
	}
	return b, nil

}

func (b Board) movePawn(start, end Coordinate, piece *pieces.Piece) (Board, error) {
	//validate that pawn can move
	if piece.MoveCount == 0 {
		//pawn can move up to 2 squares vertically
		if piece.Team == pieces.White {
			if start.X < end.X {
				return b, errors.New("white pawn must move forward")
			}

		} else if piece.Team == pieces.Black {
			if start.X > end.X {
				return b, errors.New("black pawn must move forward")
			}

		}

	}

	//passed validation move the pawn
	piece.MoveCount += 1
	fmt.Println("adding 1 to Pawn Move count")
	b.Squares[start.X][start.Y] = pieces.Piece{Type: pieces.Empty, Team: pieces.Neutral}
	b.Squares[end.X][end.Y] = *piece
	return b, nil
}

func (b Board) moveRook(start, end Coordinate, piece pieces.Piece) (Board, error) {
	b.Squares[start.X][start.Y] = pieces.Piece{Type: pieces.Empty, Team: pieces.Neutral}
	b.Squares[end.X][end.Y] = piece
	return b, nil
}

func (b Board) moveKnight(start, end Coordinate, piece pieces.Piece) (Board, error) {
	b.Squares[start.X][start.Y] = pieces.Piece{Type: pieces.Empty, Team: pieces.Neutral}
	b.Squares[end.X][end.Y] = piece
	return b, nil
}

func (b Board) moveBishop(start, end Coordinate, piece pieces.Piece) (Board, error) {
	b.Squares[start.X][start.Y] = pieces.Piece{Type: pieces.Empty, Team: pieces.Neutral}
	b.Squares[end.X][end.Y] = piece
	return b, nil
}

func (b Board) moveQueen(start, end Coordinate, piece pieces.Piece) (Board, error) {
	b.Squares[start.X][start.Y] = pieces.Piece{Type: pieces.Empty, Team: pieces.Neutral}
	b.Squares[end.X][end.Y] = piece
	return b, nil
}

func (b Board) moveKing(start, end Coordinate, piece pieces.Piece) (Board, error) {
	b.Squares[start.X][start.Y] = pieces.Piece{Type: pieces.Empty, Team: pieces.Neutral}
	b.Squares[end.X][end.Y] = piece
	return b, nil
}

func (b Board) getPiece(coord Coordinate) (pieces.Piece, error) {
	if coord.X > 7 || coord.Y > 7 || coord.X < 0 || coord.Y < 0 {
		noPiece := pieces.Piece{
			Type: 99,
			Team: 99,
		}

		return noPiece, errors.New("coordinate must be within bounds of the board")
	}
	return b.Squares[coord.X][coord.Y], nil

}

func CreateDefaultBoard() Board {
	var board Board

	// Pawns
	for file := 0; file < 8; file++ {
		board.Squares[6][file] = pieces.Piece{Type: pieces.Pawn, Team: pieces.White}
		board.Squares[5][file] = pieces.Piece{Type: pieces.Empty, Team: pieces.Neutral}
		board.Squares[4][file] = pieces.Piece{Type: pieces.Empty, Team: pieces.Neutral}
		board.Squares[3][file] = pieces.Piece{Type: pieces.Empty, Team: pieces.Neutral}
		board.Squares[2][file] = pieces.Piece{Type: pieces.Empty, Team: pieces.Neutral}
		board.Squares[1][file] = pieces.Piece{Type: pieces.Pawn, Team: pieces.Black}
	}

	// Back ranks
	backRank := []pieces.PieceType{
		pieces.Rook, pieces.Knight, pieces.Bishop, pieces.Queen,
		pieces.King, pieces.Bishop, pieces.Knight, pieces.Rook,
	}

	for file, pt := range backRank {
		board.Squares[7][file] = pieces.Piece{Type: pt, Team: pieces.White}
		board.Squares[0][file] = pieces.Piece{Type: pt, Team: pieces.Black}
	}

	return board
}
