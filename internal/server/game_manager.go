package server

import "chess/internal/chess"

type GameManager struct {
	clients map[*Client]bool
	board   *chess.Board
}

func NewGameManager() *GameManager {

	return &GameManager{
		clients: make(map[*Client]bool),
		board:   chess.NewBoard(),
	}
}

func (gm *GameManager) Register(c *Client) {
	gm.clients[c] = true
}
