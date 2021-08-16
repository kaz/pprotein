package mux

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewDebugHandler() http.Handler {
	r := mux.NewRouter()
	EnableDebugHandler(r)
	return r
}
