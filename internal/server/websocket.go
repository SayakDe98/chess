package server

import (
	"log"
	"net/http"

	"chess/internal/chess"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (gm *GameManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte),
	}

	gm.Register(client)

	go client.readPump(gm)
	go client.writePump()
}

func (c *Client) readPump(gm *GameManager) {

	for {

		var move chess.Move

		err := c.conn.ReadJSON(&move)
		if err != nil {
			break
		}

		if gm.board.IsValidMove(move) {

			gm.board.MakeMove(move)

			// state := gm.board.Serialize()
			state := gm.board.ToFEN()
			for client := range gm.clients {

				client.conn.WriteJSON(state)
			}
		}
	}
}

func (c *Client) writePump() {

	for msg := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}
