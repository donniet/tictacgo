package tictacgo

type Square uint8

const (
	Empty Square = 0
	X     Square = 1
	O     Square = 2
)

func (s Square) String() string {
	switch s {
	case X:
		return "X"
	case O:
		return "O"
	default:
		return " "
	}
}

type TicTacError string

func (e TicTacError) Error() string {
	return string(e)
}

const (
	ErrorOutOfBounds   = TicTacError("out of bounds")
	ErrorInvalidFormat = TicTacError("invalid position format")
)

type Position struct {
	pos [9]Square
}

func (p Position) String() string {
	s := ""
	for i := 0; i < 9; i++ {
		s += p.pos[i].String()
		if i%3 == 2 {
			s += "\n"
		}
	}
	return s
}

func FromString(pos string) (Position, error) {
	q := Position{}
	x, y := 0, 0
	for _, c := range pos {
		if x == 3 {
			if c != '\n' {
				return q, ErrorInvalidFormat
			}
			x = 0
			y++
			continue
		}
		if y == 3 {
			return q, ErrorInvalidFormat
		}
		switch c {
		case 'X':
			q.Set(x, y, X)
		case 'O':
			q.Set(x, y, O)
		case ' ':
			q.Set(x, y, Empty)
		default:
			return q, ErrorInvalidFormat
		}
		x++
	}
	if x != 0 || y != 3 {
		return q, ErrorInvalidFormat
	}
	return q, nil
}

func (p *Position) Set(x, y int, s Square) {
	if x < 0 || x >= 3 || y < 0 || y >= 3 {
		return
	}

	p.pos[y*3+x] = s
}

func (p Position) Get(x, y int) Square {
	if x < 0 || x >= 3 || y < 0 || y >= 3 {
		return Empty
	}

	return p.pos[y*3+x]
}

func (p *Position) Clear() {
	for i := range p.pos {
		p.pos[i] = Empty
	}
}

func (p Position) FlipX() Position {
	q := Position{}

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			q.Set(2-x, y, p.Get(x, y))
		}
	}

	return q
}

func (p Position) FlipY() Position {
	q := Position{}

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			q.Set(x, 2-y, p.Get(x, y))
		}
	}

	return q
}

func (p Position) Rotate(i int) Position {
	q := Position{}

	if i < -3 || i > 3 {
		i %= 4
	}
	if i < 0 {
		i += 4
	}

	switch i {
	case 0:
		// copy
		q.pos = p.pos
	case 1:
		// rotate once clockwise
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				q.Set(x, y, p.Get(y, 2-x))
			}
		}
	case 2:
		// rotate twice clockwise
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				q.Set(x, y, p.Get(2-x, 2-y))
			}
		}
	case 3:
		// rotate once clockwise
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				q.Set(x, y, p.Get(2-y, x))
			}
		}
	}

	return q
}

func (p Position) Turn() Square {
	countX, countO := 0, 0
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			switch p.Get(x, y) {
			case X:
				countX++
			case O:
				countO++
			}

		}
	}
	if countX <= countO {
		return X
	}
	return O
}

func (p Position) Next() []Position {
	s := p.Turn()
	var ret []Position

	// if it's a winner there are no more valid next positions
	if p.IsWin() != Empty {
		return ret
	}

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if p.Get(x, y) == Empty {
				q := p
				q.Set(x, y, s)
				ret = append(ret, q)
			}
		}
	}

	return ret
}

func (p Position) IsWin() Square {
	// vertical
	for x := 0; x < 3; x++ {
		s := p.Get(x, 0)
		if s == Empty {
			continue
		}
		if s == p.Get(x, 1) && s == p.Get(x, 2) {
			return s
		}
	}

	// horizontal
	for y := 0; y < 3; y++ {
		s := p.Get(0, y)
		if s == Empty {
			continue
		}
		if s == p.Get(1, y) && s == p.Get(2, y) {
			return s
		}
	}

	// diagonals
	if s := p.Get(1, 1); s != Empty {
		if s == p.Get(0, 0) && s == p.Get(2, 2) {
			return s
		}
		if s == p.Get(0, 2) && s == p.Get(2, 0) {
			return s
		}
	}

	return Empty
}
