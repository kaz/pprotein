package integration

import (
	"github.com/kaz/pprotein/internal/git"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"github.com/felixge/fgprof"
	"github.com/gorilla/mux"
	"github.com/kaz/pprotein/internal/tail"
)

func NewDebugHandler() http.Handler {
	r := mux.NewRouter()
	RegisterDebugHandlers(r)
	return r
}

func RegisterDebugHandlers(r *mux.Router) {
	r.Use(gitRevisionResponseMiddleware)
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
		revision := git.GetCommitHash(getEnvOrDefault("GIT_REPO_DIR", "/home/isucon/repo"))
		if revision != "" {
			rw.Header().Set("X-GIT-REVISION", strings.ReplaceAll(string(revision), "\n", ""))
		}
		next.ServeHTTP(rw, r)
	})
}

func getEnvOrDefault(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
