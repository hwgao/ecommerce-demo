package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Product service starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
