package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Order service starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
