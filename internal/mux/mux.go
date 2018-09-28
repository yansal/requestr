package mux

import (
	"html/template"
	"log"
	"net/http"

	"github.com/yansal/requestr/internal/broker"
	"github.com/yansal/requestr/internal/store"
)

func New(store *store.Store, publisher broker.Publisher, template *template.Template) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.Handle("/", rootHandler(store, publisher, template))
	return mux
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func handleError(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log.Printf("%+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
