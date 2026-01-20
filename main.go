package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Response JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"message": "API Running",
		})
	})

	fmt.Println("Starting server on port 8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("error starting server:", err)
	}
}
