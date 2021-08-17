package mux

import (
	"net/http"
	"net/http/pprof"

	"github.com/felixge/fgprof"
	"github.com/gorilla/mux"
	"github.com/kaz/pprotein/internal/tail"
)

func NewDebugHandler() http.Handler {
	r := mux.NewRouter()
	EnableDebugHandler(r)
	return r
}

func registerDebugHandlers(r *mux.Router) {
	r.Handle("/debug/log/test", tail.NewTailHandler("./test"))
	r.Handle("/debug/log/nginx", tail.NewTailHandler("/var/log/nginx/access.log"))
	r.Handle("/debug/log/mysql", tail.NewTailHandler("/var/log/mysql/mysql-slow.log"))

	r.Handle("/debug/fgprof", fgprof.Handler())

	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.HandleFunc("/debug/pprof/{h:.*}", pprof.Index)
}
