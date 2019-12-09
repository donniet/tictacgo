package tictacgo

import "io"

// AllValidPositionEquiv gets all valid tic tac toe positions under the equivalence operations of flipping and rotating
func AllValidPositionsEquiv() []Position {
	visited := make(map[Position]bool)
	stack := []Position{Position{}}

	var p Position
	for len(stack) > 0 {
		// pop off the stack
		p, stack = stack[len(stack)-1], stack[:len(stack)-1]

		next := p.Next()
		for _, q := range next {
			if !visited[q] && !visited[q.FlipX()] && !visited[q.FlipY()] && !visited[q.Rotate(1)] && !visited[q.Rotate(2)] && !visited[q.Rotate(3)] {
				visited[q] = true
				stack = append(stack, q)
			}
		}
	}

	var ret []Position
	for p := range visited {
		ret = append(ret, p)
	}
	return ret
}

// ReadPositions from a reader and return slice of positions or error
func ReadPositions(r io.Reader) ([]Position, error) {
	var ret []Position

	// read 12 bytes
	buf := make([]byte, 12)

	for {
		if _, err := io.ReadFull(r, buf); err == io.EOF {
			break
		} else if err != nil {
			return ret, err
		}

		if p, err := FromString(string(buf)); err != nil {
			return ret, err
		} else {
			ret = append(ret, p)
		}
	}

	return ret, nil
}
