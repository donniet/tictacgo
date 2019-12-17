package tictacgo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
)

type Scorer interface {
	Score(Position, Square) float64
}

type Tally map[Square]int // tally of wins or cats games

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
	positions map[Position]Tally
	Curiosity float64 // value from 0 to 1
}

type jsonEvaluation struct {
	Positions map[Position]Tally `json:"positions"`
	Curiosity float64            `json:"curiosity"`
}

func NewEvaluation() *Evaluation {
	return &Evaluation{positions: make(map[Position]Tally)}
}

func (e *Evaluation) Load(r io.Reader) error {
	var t jsonEvaluation

	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	e.positions = t.Positions
	e.Curiosity = t.Curiosity

	return nil
}

func (e *Evaluation) Save(w io.Writer) error {
	t := jsonEvaluation{
		Positions: e.positions,
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
	for _, m := range moves {
		sym := m.Symmetries()
		found := false
		for _, p := range sym {
			if t, ok := e.positions[p]; ok {
				t[winFor]++
				found = true
				break
			}
		}

		if !found {
			e.positions[m] = make(Tally)
			e.positions[m][winFor]++
		}
	}
}

func (e *Evaluation) getTallyWithSymmetry(pos Position) Tally {
	sym := pos.Symmetries()

	for _, p := range sym {
		if t, ok := e.positions[p]; ok {
			return t
		}
	}

	return Tally{}
}

func (e *Evaluation) Score(pos Position, turn Square) float64 {
	t := e.getTallyWithSymmetry(pos)

	unadjusted := float64(t[Empty])

	switch turn {
	case X:
		unadjusted += float64(t[X])
		unadjusted -= float64(t[O])
	case O:
		unadjusted -= float64(t[X])
		unadjusted += float64(t[O])
	}

	// adjust it by curiosity
	return unadjusted + e.Curiosity*rand.Float64()
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
