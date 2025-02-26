package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/NikolaiMarkalainen/Router/api"
	"github.com/NikolaiMarkalainen/Router/utils"
)

func main() {
	// Initialize the router
	router := api.NewRouter()

  router.GET("/", func(w *utils.ResponseWriter, r *http.Request) {
		fmt.Print("home")
  })
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		return
	}

	// Serve the requests with the router
	if err := http.Serve(l, router); err != nil {
		fmt.Printf("Server closed: %s\n", err)
	}
}
