package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/donniet/tictacgo"
)

func ParsePosition(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pos := r.URL.Query().Get("p")

		var p tictacgo.Position
		var err error

		if pos != "" {
			p, err = tictacgo.FromString(pos)

			if err != nil {
				http.Error(w, fmt.Sprintf("invalid position: '%s'", pos), http.StatusBadRequest)
				return
			}
		}

		if r.Method == http.MethodPost {
			// we should have a list of positions and a new move to pull from the form
			sh, _ := url.QueryUnescape(r.FormValue("history"))
			sx := r.FormValue("row")
			sy := r.FormValue("col")

			var history []tictacgo.Position
			if err := json.Unmarshal([]byte(sh), &history); err != nil {
				http.Error(w, fmt.Sprintf("invalid history field: '%s'", sh), http.StatusBadRequest)
				return
			}

			var err error
			var x, y int

			if x, err = strconv.Atoi(sx); err != nil || x < 0 || x > 2 {
				http.Error(w, fmt.Sprintf("x value is error: '%s'", sx), http.StatusBadRequest)
				return
			}
			if y, err = strconv.Atoi(sy); err != nil || y < 0 || y > 2 {
				http.Error(w, fmt.Sprintf("y value error: %s", sy), http.StatusBadRequest)
				return
			}

			var cur tictacgo.Position
			if len(history) > 0 {
				cur = history[0]
			}

			if cur.Get(x, y) != tictacgo.Empty || cur.IsWin() != tictacgo.Empty {
				http.Error(w, fmt.Sprintf("invalid move: %d,%d", x, y), http.StatusBadRequest)
				return
			}

			cur.Set(x, y, cur.Turn())
			history = append([]tictacgo.Position{cur}, history...)

			r = r.WithContext(context.WithValue(r.Context(), ContextHistoryKey, history))
			r = r.WithContext(context.WithValue(r.Context(), ContextPositionKey, &cur))
		} else if r.Method == http.MethodGet {
			history := []tictacgo.Position{p}
			r = r.WithContext(context.WithValue(r.Context(), ContextHistoryKey, history))
			r = r.WithContext(context.WithValue(r.Context(), ContextPositionKey, &p))
		}

		handler.ServeHTTP(w, r)
	})
}

func ParsePositionFunc(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return ParsePosition(http.HandlerFunc(handler))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mux := http.NewServeMux()

	server := &Server{
		Eval: tictacgo.NewEvaluation(),
	}

	mux.Handle("/game", ParsePosition(server))
	mux.HandleFunc("/positions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := server.Eval.Save(w); err != nil {
			log.Printf("error saving evaluations: %s", err)
		}
	})
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		server.Eval.Reset()

		http.Redirect(w, r, "/game", http.StatusTemporaryRedirect)
	})

	http.ListenAndServe(":8080", mux)
}
