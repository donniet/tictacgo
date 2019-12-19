package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/donniet/tictacgo"
)

var (
	positionTemplate *template.Template
)

type ServerContextKey int

const (
	ContextPositionKey ServerContextKey = iota
	ContextHistoryKey
)

func init() {
	positionTemplate = template.Must(template.New("position").Funcs(map[string]interface{}{
		"X":           func() tictacgo.Square { return tictacgo.X },
		"O":           func() tictacgo.Square { return tictacgo.O },
		"Empty":       func() tictacgo.Square { return tictacgo.Empty },
		"mul":         func(x, y int) int { return x * y },
		"postAction":  func() string { return "game" },
		"queryEscape": func(str string) string { return url.QueryEscape(str) },
	}).ParseGlob("templates/*.html"))
}

// WritePosition writes out the position using the embedded template to a Writer
func WritePosition(w io.Writer, p tictacgo.Position) {
	positionTemplate.ExecuteTemplate(w, "position.html", p)
}

type Server struct {
	Eval *tictacgo.Evaluation
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := r.Context()
	position, ok := context.Value(ContextPositionKey).(*tictacgo.Position)
	if !ok {
		http.Error(w, "error getting position from context", http.StatusInternalServerError)
		return
	}

	history, ok := context.Value(ContextHistoryKey).([]tictacgo.Position)
	if !ok {
		http.Error(w, "error getting history from context", http.StatusInternalServerError)
		return
	}

	if position.IsComplete() {
		// update the history
		s.Eval.Result(history, position.IsWin())
	} else if len(history) > 1 || rand.Intn(2) == 0 {
		// if this isn't the first move, or we flipped a coin and won
		pos, err := s.Eval.ChooseNext(*position)

		if err != nil {
			http.Error(w, "error choosing next move", http.StatusInternalServerError)
			log.Printf("error chosing next position: %s", err)
			return
		}

		position = &pos
		history = append([]tictacgo.Position{pos}, history...)

		if position.IsComplete() {
			s.Eval.Result(history, position.IsWin())
		}
	}

	b, err := json.Marshal(history)
	if err != nil {
		http.Error(w, "error marshalling history to json", http.StatusInternalServerError)
		log.Printf("error marshalling history: %s", err)
		return
	}

	if err := positionTemplate.ExecuteTemplate(w, "game.html", struct {
		Position    *tictacgo.Position
		History     []tictacgo.Position
		HistoryJSON string
	}{
		position,
		history,
		string(b),
	}); err != nil {
		log.Printf("error executing template: %v", err)
	}
}
