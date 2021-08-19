// +build !release

package standalone

import (
	"log"
	"net/http"
	"os"

	"github.com/kaz/pprotein/integration/mux"
)

func Integrate() {
	port := os.Getenv("PORT_DEBUG")
	if port == "" {
		port = "8000"
	}

	log.Println("Launching debug server ...")
	go func() {
		if err := http.ListenAndServe(":"+port, mux.NewDebugHandler()); err != nil {
			log.Printf("failed to start debug server: %v\n", err)
		}
	}()
}
