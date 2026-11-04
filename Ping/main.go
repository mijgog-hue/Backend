package main

import (
	"fmt"
	"log"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func main() {
	http.HandleFunc("/ping", pingHandler)
	log.Println("server started on: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
