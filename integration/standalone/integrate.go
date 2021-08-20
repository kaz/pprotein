package standalone

import (
	"log"
	"net/http"

	"github.com/kaz/pprotein/integration"
)

func Integrate(addr string) {
	log.Printf("[DEBUG_SERVER] Listening on %v\n", addr)
	if err := http.ListenAndServe(addr, integration.NewDebugHandler()); err != nil {
		log.Printf("failed to start debug server: %v\n", err)
	}
}
