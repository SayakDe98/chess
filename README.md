# Chess Multiplayer Server (Go)

A real-time multiplayer chess server written in **Go** using **WebSockets**.

This backend handles:

* chess rules
* move validation
* check / checkmate detection
* castling and en passant
* FEN board format
* multiplayer synchronization via WebSockets

The frontend connects using a WebSocket client.

---

# Features

* Full chess rule engine
* Check detection
* Checkmate detection
* Castling
* En passant
* FEN board representation
* WebSocket multiplayer updates

---

# Project Structure

```
chess-game/

├── cmd/
│   └── server/
│        main.go
│
├── internal/
│   ├── chess/
│   │    board.go
│   │    piece.go
│   │    move.go
│   │    move_generator.go
│   │    pawn.go
│   │    rook.go
│   │    knight.go
│   │    bishop.go
│   │    queen.go
│   │    king.go
│   │
│   └── server/
│        websocket.go
│        game_manager.go
│
└── go.mod
```

---

# Requirements

* Go 1.22+
* Git

---

# Install dependencies

```
go mod tidy
```

---

# Run the server

```
go run ./cmd/server
```

Server will start on:

```
http://localhost:8080
```

WebSocket endpoint:

```
ws://localhost:8080/ws
```

---

# WebSocket API

Clients send moves as JSON:

```
{
  "fromX": 4,
  "fromY": 1,
  "toX": 4,
  "toY": 3
}
```

Server broadcasts updated board state.

Example response:

```
{
  "fen": "rnbqkbnr/pppppppp/8/8/4p3/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
}
```

---

# Chess Engine

The engine supports:

* legal move validation
* check detection
* checkmate detection
* castling
* en passant
* FEN serialization

---

# Future Improvements

* matchmaking
* player authentication
* game rooms
* PGN recording
* chess clock
* spectator mode

---

# License

MIT
