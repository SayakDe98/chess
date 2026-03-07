package chess

type Board struct {
	Squares [8][8]*Piece
	Turn    Color

	WhiteKingMoved bool
	BlackKingMoved bool

	WhiteRookAMoved bool
	WhiteRookHMoved bool
	BlackRookAMoved bool
	BlackRookHMoved bool

	EnPassantX int
	EnPassantY int
}

func NewBoard() *Board {
	b := &Board{}

	b.EnPassantX = -1
	b.EnPassantY = -1

	b.setup()

	return b
}

func (b *Board) setup() {

	for i := 0; i < 8; i++ {
		b.Squares[1][i] = &Piece{Pawn, White}
		b.Squares[6][i] = &Piece{Pawn, Black}
	}

	b.Squares[0][0] = &Piece{Rook, White}
	b.Squares[0][7] = &Piece{Rook, White}
	b.Squares[7][0] = &Piece{Rook, Black}
	b.Squares[7][7] = &Piece{Rook, Black}

	b.Squares[0][1] = &Piece{Knight, White}
	b.Squares[0][6] = &Piece{Knight, White}
	b.Squares[7][1] = &Piece{Knight, Black}
	b.Squares[7][6] = &Piece{Knight, Black}

	b.Squares[0][2] = &Piece{Bishop, White}
	b.Squares[0][5] = &Piece{Bishop, White}
	b.Squares[7][2] = &Piece{Bishop, Black}
	b.Squares[7][5] = &Piece{Bishop, Black}

	b.Squares[0][3] = &Piece{Queen, White}
	b.Squares[7][3] = &Piece{Queen, Black}

	b.Squares[0][4] = &Piece{King, White}
	b.Squares[7][4] = &Piece{King, Black}

	b.Turn = White
}

func (b *Board) MakeMove(m Move) {

	p := b.Squares[m.FromY][m.FromX]

	// reset en passant square
	b.EnPassantX = -1
	b.EnPassantY = -1

	// pawn double move (enable en passant)
	if p.Type == Pawn && abs(m.ToY-m.FromY) == 2 {
		b.EnPassantX = m.FromX
		b.EnPassantY = (m.FromY + m.ToY) / 2
	}

	// en passant capture
	if p.Type == Pawn && m.ToX == b.EnPassantX && m.ToY == b.EnPassantY {

		if p.Color == White {
			b.Squares[m.ToY-1][m.ToX] = nil
		} else {
			b.Squares[m.ToY+1][m.ToX] = nil
		}
	}

	// move the piece
	b.Squares[m.ToY][m.ToX] = p
	b.Squares[m.FromY][m.FromX] = nil

	// castling rook movement
	if p.Type == King && abs(m.ToX-m.FromX) == 2 {

		// king side
		if m.ToX == 6 {
			rook := b.Squares[m.FromY][7]
			b.Squares[m.FromY][5] = rook
			b.Squares[m.FromY][7] = nil
		}

		// queen side
		if m.ToX == 2 {
			rook := b.Squares[m.FromY][0]
			b.Squares[m.FromY][3] = rook
			b.Squares[m.FromY][0] = nil
		}
	}

	// update castling rights
	if p.Type == King {

		if p.Color == White {
			b.WhiteKingMoved = true
		} else {
			b.BlackKingMoved = true
		}
	}

	if p.Type == Rook {

		if p.Color == White {
			if m.FromX == 0 && m.FromY == 0 {
				b.WhiteRookAMoved = true
			}
			if m.FromX == 7 && m.FromY == 0 {
				b.WhiteRookHMoved = true
			}
		}

		if p.Color == Black {
			if m.FromX == 0 && m.FromY == 7 {
				b.BlackRookAMoved = true
			}
			if m.FromX == 7 && m.FromY == 7 {
				b.BlackRookHMoved = true
			}
		}
	}

	// switch turn
	if b.Turn == White {
		b.Turn = Black
	} else {
		b.Turn = White
	}
}

func (b *Board) Clone() *Board {

	nb := &Board{}
	nb.Turn = b.Turn

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {

			if b.Squares[y][x] != nil {

				p := *b.Squares[y][x]
				nb.Squares[y][x] = &p
			}
		}
	}

	return nb
}

func (b *Board) IsValidMove(m Move) bool {
	p := b.Squares[m.FromY][m.FromX]

	if p == nil {
		return false
	}

	if p.Color != b.Turn {
		return false
	}

	moves := b.GenerateMoves(m.FromX, m.FromY)

	for _, mv := range moves {

		if mv.ToX == m.ToX && mv.ToY == m.ToY {

			test := b.Clone()
			test.MakeMove(m)

			if test.InCheck(p.Color) {
				return false
			}

			return true
		}
	}

	return false
}

func (b *Board) InCheck(color Color) bool {

	var kingX, kingY int

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {

			p := b.Squares[y][x]

			if p != nil && p.Type == King && p.Color == color {
				kingX = x
				kingY = y
			}
		}
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {

			p := b.Squares[y][x]

			if p == nil || p.Color == color {
				continue
			}

			moves := b.GenerateMoves(x, y)

			for _, m := range moves {
				if m.ToX == kingX && m.ToY == kingY {
					return true
				}
			}
		}
	}

	return false
}

func (b *Board) Serialize() [][]string {

	board := make([][]string, 8)

	for y := 0; y < 8; y++ {

		row := make([]string, 8)

		for x := 0; x < 8; x++ {

			p := b.Squares[y][x]

			if p == nil {
				row[x] = "."
				continue
			}

			row[x] = pieceSymbol(p)
		}

		board[y] = row
	}

	return board
}

func pieceSymbol(p *Piece) string {

	symbols := map[PieceType]string{
		Pawn:   "p",
		Rook:   "r",
		Knight: "n",
		Bishop: "b",
		Queen:  "q",
		King:   "k",
	}

	s := symbols[p.Type]

	if p.Color == White {
		return string(s[0] - 32)
	}

	return s
}

func (b *Board) IsCheckmate(color Color) bool {

	if !b.InCheck(color) {
		return false
	}

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {

			p := b.Squares[y][x]

			if p == nil || p.Color != color {
				continue
			}

			moves := b.GenerateMoves(x, y)

			for _, m := range moves {

				test := b.Clone()

				if test.IsValidMove(m) {
					return false
				}
			}
		}
	}

	return true
}

func (b *Board) ToFEN() string {

	fen := ""

	for y := 7; y >= 0; y-- {

		empty := 0

		for x := 0; x < 8; x++ {

			p := b.Squares[y][x]

			if p == nil {
				empty++
				continue
			}

			if empty > 0 {
				fen += string(rune(empty + '0'))
				empty = 0
			}

			fen += pieceSymbol(p)
		}

		if empty > 0 {
			fen += string(rune(empty + '0'))
		}

		if y > 0 {
			fen += "/"
		}
	}

	// turn
	if b.Turn == White {
		fen += " w "
	} else {
		fen += " b "
	}

	// castling rights
	castle := ""

	if !b.WhiteKingMoved && !b.WhiteRookHMoved {
		castle += "K"
	}
	if !b.WhiteKingMoved && !b.WhiteRookAMoved {
		castle += "Q"
	}
	if !b.BlackKingMoved && !b.BlackRookHMoved {
		castle += "k"
	}
	if !b.BlackKingMoved && !b.BlackRookAMoved {
		castle += "q"
	}

	if castle == "" {
		castle = "-"
	}

	fen += castle + " "

	// en passant
	if b.EnPassantX == -1 {
		fen += "-"
	} else {
		file := string(rune('a' + b.EnPassantX))
		rank := string(rune('1' + b.EnPassantY))

		fen += file + rank
	}

	return fen
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
