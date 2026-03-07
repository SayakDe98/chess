package chess

func KnightMoves(b *Board, x, y int) []Move {

	var moves []Move

	offsets := [][]int{
		{2, 1}, {1, 2},
		{-1, 2}, {-2, 1},
		{-2, -1}, {-1, -2},
		{1, -2}, {2, -1},
	}

	for _, o := range offsets {

		nx := x + o[0]
		ny := y + o[1]

		if nx >= 0 && nx < 8 && ny >= 0 && ny < 8 {
			moves = append(moves, Move{x, y, nx, ny})
		}
	}

	return moves
}
