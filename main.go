package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}
func byeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Bye!")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/bye", byeHandler)
	log.Println("server started on: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
