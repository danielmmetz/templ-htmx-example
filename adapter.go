package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

type HandlerFunc func(*http.Request) (templ.Component, error)

func Handle(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		component, err := h(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "error:", err)
			return
		}
		if err := component.Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "error rendering component:", err)
			return
		}
	}
}
