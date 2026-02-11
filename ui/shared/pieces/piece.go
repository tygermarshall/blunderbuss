package pieces

type PieceType int

type Team int

func (team Team) String() string {
	switch team {
	case White:
		return "white"
	case Black:
		return "black"
	case Neutral:
		return "neutral"
	default:
		return "unknown team"
	}
}

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
	Empty
)

func (piece PieceType) String() string {
	switch piece {
	case Pawn:
		return "pawn"
	case Knight:
		return "knight"
	case Bishop:
		return "bishop"
	case Rook:
		return "rook"
	case Queen:
		return "queen"
	case King:
		return "king"
	default:
		return "unknown piece"
	}
}

const (
	White Team = iota
	Black
	Neutral
)

type Piece struct {
	Type      PieceType
	Team      Team
	MoveCount int
}
