package main

import (
	"os"

	"github.com/kaz/pprotein/integration/standalone"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "19000"
	}
	standalone.Integrate(":" + port)
}
