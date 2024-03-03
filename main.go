package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {

		param := chi.URLParam(r, "id")

		render.JSON(w, r, param)

	})

	http.ListenAndServe(":3002", r)
}
