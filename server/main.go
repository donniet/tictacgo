package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

		r = r.WithContext(context.WithValue(context.Background(), ContextPositionKey, &p))

		handler.ServeHTTP(w, r)
	})
}

func ParsePositionFunc(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return ParsePosition(http.HandlerFunc(handler))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mux := http.NewServeMux()

	mux.Handle("/game", ParsePositionFunc(ServePosition))

	http.ListenAndServe(":8080", mux)
}
