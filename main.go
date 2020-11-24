package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	h := http.HandlerFunc(Echo)

	log.Println("Listening and serving on localhost:8000")
	if err := http.ListenAndServe(":8000", h); err != nil {
		log.Fatal(err)
	}
}

// Echo just tells you the request you made
func Echo(w http.ResponseWriter, r *http.Request) {
	id := rand.Intn(1000)
	fmt.Println("Starting", id)
	time.Sleep(3 * time.Second)
	fmt.Fprintln(w, "You asked to", r.Method, r.URL.Path)
	fmt.Println("ending", id)
}
