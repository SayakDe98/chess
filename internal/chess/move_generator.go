package chess

func (b *Board) GenerateMoves(x, y int) []Move {

	p := b.Squares[y][x]

	if p == nil {
		return nil
	}

	switch p.Type {

	case Pawn:
		return PawnMoves(b, x, y)

	case Rook:
		return RookMoves(b, x, y)

	case Knight:
		return KnightMoves(b, x, y)

	case Bishop:
		return BishopMoves(b, x, y)

	case Queen:
		return QueenMoves(b, x, y)

	case King:
		return KingMoves(b, x, y)
	}

	return nil
}
