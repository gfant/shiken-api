package main

import (
	"executor/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/runCode", handler.RunCode)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
