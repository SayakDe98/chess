package chess

func RookMoves(b *Board, x, y int) []Move {

	var moves []Move

	dirs := [][]int{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	}

	for _, d := range dirs {

		nx := x
		ny := y

		for {

			nx += d[0]
			ny += d[1]

			if nx < 0 || nx >= 8 || ny < 0 || ny >= 8 {
				break
			}

			target := b.Squares[ny][nx]

			if target == nil {

				moves = append(moves, Move{x, y, nx, ny})
				continue
			}

			if target.Color != b.Squares[y][x].Color {
				moves = append(moves, Move{x, y, nx, ny})
			}

			break
		}
	}

	return moves
}
