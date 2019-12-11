package main

import (
	"html/template"
	"io"
	"log"
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
)

func init() {
	positionTemplate = template.Must(template.New("position").Funcs(map[string]interface{}{
		"X":           func() tictacgo.Square { return tictacgo.X },
		"O":           func() tictacgo.Square { return tictacgo.O },
		"Empty":       func() tictacgo.Square { return tictacgo.Empty },
		"mul":         func(x, y int) int { return x * y },
		"postAction":  func() string { return "game" },
		"queryEscape": func(str string) string { return url.QueryEscape(str) },
	}).ParseGlob("../templates/*.html"))
}

// WritePosition writes out the position using the embedded template to a Writer
func WritePosition(w io.Writer, p tictacgo.Position) {
	positionTemplate.ExecuteTemplate(w, "position.html", p)
}

func ServePosition(w http.ResponseWriter, r *http.Request) {
	context := r.Context()
	position, ok := context.Value(ContextPositionKey).(*tictacgo.Position)

	if !ok {
		http.Error(w, "error getting position from context", http.StatusInternalServerError)
		return
	}

	if err := positionTemplate.ExecuteTemplate(w, "game.html", position); err != nil {
		log.Printf("error executing template: %v", err)
	}
}
