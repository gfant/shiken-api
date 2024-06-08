package main

import (
	"fmt"
	"net/http"
	"os"
	"shiken_server/handler"
)

func main() {
	err := os.MkdirAll("tmp/pro/pre", 0755)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/getCode", handler.GetCode)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
