package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/tygermarshall/blunderbuss/shared"
	"github.com/tygermarshall/blunderbuss/shared/board"
)

var ErrGameNotFound = errors.New("game not found")

type Game struct {
	Board      board.Board
	TurnNumber int
}

// gameEntry holds a game and its own mutex so one game's operations
// don't block others.
type gameEntry struct {
	mu   sync.Mutex
	game *Game
}

// GameStore holds all games. Use the map mutex for create/lookup;
// use each gameEntry's mutex for reading or mutating that game.
type GameStore struct {
	mu    sync.RWMutex
	games map[string]*gameEntry
}

func NewGameStore() *GameStore {
	return &GameStore{games: make(map[string]*gameEntry)}
}

func newGame() *Game {
	return &Game{
		Board:      board.CreateDefaultBoard(),
		TurnNumber: 1,
	}
}

func generateID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Create creates a new game, stores it, and returns its ID.
func (s *GameStore) Create() (string, error) {
	id, err := generateID()
	if err != nil {
		return "", err
	}
	g := newGame()
	s.mu.Lock()
	s.games[id] = &gameEntry{game: g}
	shared.PrintBoard(&g.Board)
	s.mu.Unlock()
	return id, nil
}

// Get returns a copy of the game state for the given ID, or nil, false if not found.
func (s *GameStore) Get(id string) (*Game, bool) {
	s.mu.RLock()
	entry := s.games[id]
	s.mu.RUnlock()
	if entry == nil {
		return nil, false
	}
	entry.mu.Lock()
	cp := &Game{Board: entry.game.Board, TurnNumber: entry.game.TurnNumber}
	entry.mu.Unlock()
	return cp, true
}

// Move applies a move to the game. Returns ErrGameNotFound or the board move error.
func (s *GameStore) Move(id string, from, to board.Coordinate) error {
	s.mu.RLock()
	entry := s.games[id]
	s.mu.RUnlock()
	if entry == nil {
		return ErrGameNotFound
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	newBoard, err := entry.game.Board.MovePiece(from, to)
	if err != nil {
		return err
	}
	entry.game.Board = newBoard
	entry.game.TurnNumber++
	return nil
}
