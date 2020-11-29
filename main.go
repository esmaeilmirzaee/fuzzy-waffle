package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// App starting

	log.Printf("main: Started")
	defer log.Println("main: completed")

	// Start API Server
	api := http.Server{
		Addr:         ":8000",
		Handler:      http.HandlerFunc(Echo),
		ReadTimeout:  5 * time.Second, // TODO: Why?
		WriteTimeout: 5 * time.Second, // TODO: Why?
	}

	// Make a channel to listen for errors coming from the listener
	// Use a buffered channel so the goroutine can exit if we don't
	// collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main: API listening on %s.", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("error: Listening and serving -- %s", err)

	case <-shutdown:
		log.Println("main: Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main: Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main: Could not stop server gracefully : %v", err)
		}
	}
}

// Echo just tells you about the request you made
func Echo(w http.ResponseWriter, r *http.Request) {

	// Printing a random number at the beginning and ending of a request.
	id := rand.Intn(1000)
	log.Println("start: ", id)
	defer log.Println("end: ", id)

	// Stimulate a long running results.
	time.Sleep(5 * time.Second)

	if _, err := fmt.Fprintf(w, "You asked to %s %s.", r.Method, r.URL.Path); err != nil {
		log.Printf("Echo: response to request -- %v", err)
	}
	fmt.Println("You asked to", r.Method, r.URL.Path)
}
