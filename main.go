package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	h := http.HandlerFunc(Echo)

	log.Println("Listening and serving on localhost:8000")
	if err := http.ListenAndServe(":8000", h); h != nil {
		log.Fatal(err)
	}
}

// Echo just tells about the request you made
func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You asked to", r.Method, r.URL.Path, r.Proto, r.ProtoMajor, r.ProtoMinor)
}
