package shared

import (
	"github.com/tygermarshall/blunderbuss/shared/board"
)

type CreateGameReponse struct {
	GameId string      `json:"gameId"`
	Board  board.Board `json:"board"`
}
