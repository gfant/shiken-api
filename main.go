package main

import (
	"fmt"
	"net/http"
	"shiken_server/handler"
)

func main() {
	http.HandleFunc("/runCode", handler.RunCode)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
