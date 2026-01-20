package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting server on port 8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("error starting server:", err)
	}
}
