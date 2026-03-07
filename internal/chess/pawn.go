package chess

func PawnMoves(b *Board, x, y int) []Move {

	p := b.Squares[y][x]

	dir := 1
	startRow := 1

	if p.Color == Black {
		dir = -1
		startRow = 6
	}

	var moves []Move

	// forward
	if b.Squares[y+dir][x] == nil {
		moves = append(moves, Move{x, y, x, y + dir})

		// double move
		if y == startRow && b.Squares[y+2*dir][x] == nil {
			moves = append(moves, Move{x, y, x, y + 2*dir})
		}
	}

	// capture
	for _, dx := range []int{-1, 1} {

		nx := x + dx
		ny := y + dir

		if nx < 0 || nx >= 8 || ny < 0 || ny >= 8 {
			continue
		}

		target := b.Squares[ny][nx]

		if target != nil && target.Color != p.Color {
			moves = append(moves, Move{x, y, nx, ny})
		}

		// en passant
		if nx == b.EnPassantX && ny == b.EnPassantY {
			moves = append(moves, Move{x, y, nx, ny})
		}
	}

	return moves
}
