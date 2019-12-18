package tictacgo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"sync"
)

type Scorer interface {
	Score(Position, Square) float64
}

type EvaluationError string

func (e EvaluationError) Error() string {
	return string(e)
}
func (e EvaluationError) String() string {
	return "error: " + string(e)
}

const (
	ErrorNoValidMoves = EvaluationError("no valid next moves")
)

type Evaluation struct {
	// how good is this position for X (negative is good for O)
	positions map[Position]int
	m         sync.Locker
	Curiosity float64 // value from 0 to 1
}

type jsonEvaluation struct {
	Positions map[string]int `json:"positions"`
	Curiosity float64        `json:"curiosity"`
}

func NewEvaluation() *Evaluation {
	return &Evaluation{
		positions: make(map[Position]int),
		m:         &sync.Mutex{},
	}
}

func (e *Evaluation) Load(r io.Reader) error {
	e.m.Lock()
	defer e.m.Unlock()

	var t jsonEvaluation

	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	e.positions = make(map[Position]int)

	for k, v := range t.Positions {
		pos, err := FromString(k)
		if err != nil {
			return err
		}
		e.positions[pos] = v
	}
	e.Curiosity = t.Curiosity

	return nil
}

func (e *Evaluation) Reset() {
	e.m.Lock()
	defer e.m.Unlock()

	e.positions = make(map[Position]int)
}

func (e *Evaluation) Save(w io.Writer) error {
	e.m.Lock()
	defer e.m.Unlock()

	pos := make(map[string]int)

	for k, v := range e.positions {
		pos[k.String()] = v
	}

	t := jsonEvaluation{
		Positions: pos,
		Curiosity: e.Curiosity,
	}
	if b, err := json.Marshal(&t); err != nil {
		return err
	} else if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

func (e *Evaluation) Result(moves []Position, winFor Square) {
	e.m.Lock()
	defer e.m.Unlock()

	for _, m := range moves {
		sym := m.Symmetries()
		found := false
		for s, p := range sym {
			if _, ok := e.positions[p]; ok {
				if winFor == Empty {
					e.positions[p]-- // cats games reward O
				} else if (winFor == X && s != SymmetryFlipSymbols) || (winFor == O && s == SymmetryFlipSymbols) {
					e.positions[p] += 5 - len(moves)/2 // reward the length of the game
				} else {
					e.positions[p] -= 5 - len(moves)/2
				}
				found = true
				break
			}
		}

		if !found {
			if winFor == Empty {
				e.positions[m]--
			} else if winFor == X {
				e.positions[m] += 5 - len(moves)/2
			} else {
				e.positions[m] -= 5 - len(moves)/2
			}
		}
	}
}

// must be called under lock
func (e *Evaluation) getTallyWithSymmetry(pos Position) (int, Symmetry) {
	sym := pos.Symmetries()

	for sym, p := range sym {
		if t, ok := e.positions[p]; ok {
			return t, sym
		}
	}

	return 0, SymmetryIdentity
}

func (e *Evaluation) Score(pos Position, turn Square) float64 {
	e.m.Lock()
	defer e.m.Unlock()

	t, sym := e.getTallyWithSymmetry(pos)

	score := 0.

	if (turn == X && sym != SymmetryFlipSymbols) || (turn == O && sym == SymmetryFlipSymbols) {
		score = float64(t)
	} else {
		score = float64(-t)
	}

	// adjust it by curiosity
	return score + e.Curiosity*rand.Float64()
}

func (e *Evaluation) ChooseNext(pos Position) (Position, error) {
	var ret Position
	next := pos.Next()
	if len(next) == 0 {
		return ret, ErrorNoValidMoves
	}

	turn := pos.Turn()
	ret = next[0]
	bestScore := e.Score(ret, turn)

	for i := 1; i < len(next); i++ {
		if score := e.Score(next[i], turn); score > bestScore {
			ret = next[i]
			bestScore = score
		}
	}

	return ret, nil
}
