package chess

func KingMoves(b *Board, x, y int) []Move {

	var moves []Move

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {

			if dx == 0 && dy == 0 {
				continue
			}

			nx := x + dx
			ny := y + dy

			if nx >= 0 && nx < 8 && ny >= 0 && ny < 8 {
				moves = append(moves, Move{x, y, nx, ny})
			}
		}
	}

	// castling
	p := b.Squares[y][x]

	if p.Color == White && !b.WhiteKingMoved {

		if !b.WhiteRookHMoved &&
			b.Squares[0][5] == nil &&
			b.Squares[0][6] == nil {

			moves = append(moves, Move{x, y, 6, 0})
		}

		if !b.WhiteRookAMoved &&
			b.Squares[0][1] == nil &&
			b.Squares[0][2] == nil &&
			b.Squares[0][3] == nil {

			moves = append(moves, Move{x, y, 2, 0})
		}
	}

	if p.Color == Black && !b.BlackKingMoved {

		if !b.BlackRookHMoved &&
			b.Squares[7][5] == nil &&
			b.Squares[7][6] == nil {

			moves = append(moves, Move{x, y, 6, 7})
		}

		if !b.BlackRookAMoved &&
			b.Squares[7][1] == nil &&
			b.Squares[7][2] == nil &&
			b.Squares[7][3] == nil {

			moves = append(moves, Move{x, y, 2, 7})
		}
	}

	return moves
}
