package chess

func QueenMoves(b *Board, x, y int) []Move {

	moves := RookMoves(b, x, y)
	moves = append(moves, BishopMoves(b, x, y)...)

	return moves
}
