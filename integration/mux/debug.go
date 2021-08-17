// +build !release

package mux

import (
	"github.com/gorilla/mux"
)

func Integrate(r *mux.Router) {
	EnableDebugHandler(r)
	EnableDebugMode(r)
}

func EnableDebugHandler(r *mux.Router) {
	registerDebugHandlers(r)
}

func EnableDebugMode(r *mux.Router) {
	return
}
