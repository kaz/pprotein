package integration

import (
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/felixge/fgprof"
	"github.com/gorilla/mux"
	"github.com/kaz/pprotein/internal/tail"
)

func NewDebugHandler() http.Handler {
	r := mux.NewRouter()
	r.Use(gitRevisionResponseMiddleware)
	RegisterDebugHandlers(r)
	return r
}

func RegisterDebugHandlers(r *mux.Router) {
	r.Handle("/debug/log/httplog", tail.NewTailHandler("/var/log/nginx/access.log"))
	r.Handle("/debug/log/slowlog", tail.NewTailHandler("/var/log/mysql/mysql-slow.log"))

	r.Handle("/debug/fgprof", fgprof.Handler())

	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.HandleFunc("/debug/pprof/{h:.*}", pprof.Index)
}

func gitRevisionResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		revision, err := ioutil.ReadFile("/home/isucon/.git-revision")
		if err == nil {
			rw.Header().Set("X-GIT-REVISION", strings.ReplaceAll(string(revision), "\n", ""))
		}
		next.ServeHTTP(rw, r)
	})
}
