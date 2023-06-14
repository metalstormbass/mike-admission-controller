package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/metalstormbass/mike-admission-controller/src/router"
	"github.com/rs/zerolog/log"
)

// Variable
var port string = os.Getenv("PORT")

func main() {
	// Define Router
	r := router.Router()
	port_connect := fmt.Sprintf(":%s", port)

	// Read in Certs

	// Define server

	log.Printf("HTTP Server listening on port %s", port)
	log.Print(http.ListenAndServeTLS(port_connect, "server.crt",
		"server.key", r))

}
