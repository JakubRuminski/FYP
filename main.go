package main

import (
	"fmt"
	"net/http"

	"github.com/jakubruminski/FYP/go/router/mux"
)


func main() {
	port, mux, ok := mux.INIT()

	if !ok { fmt.Println("Failed to initialize router"); return }

	fmt.Println("\033[H\033[2J")
	fmt.Println("Listening on port http://localhost:" + port)

	http.ListenAndServe(":"+port, mux)
}