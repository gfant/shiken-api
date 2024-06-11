package main

import (
	"fmt"
	"net/http"
	"reader/handler"
)

func main() {
	http.HandleFunc("/getProblemList", handler.GetProblemList)
	http.HandleFunc("/getProblem", handler.GetProblem)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
