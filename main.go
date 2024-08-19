package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// In-memory store for Items
	items := make(map[int]*Item)
	itemIDCounter := 1
	var itemMU sync.Mutex

	mux := http.NewServeMux()

	mux.HandleFunc("POST /items", CreateItemHandler(&items, &itemIDCounter, &itemMU))
	mux.HandleFunc("GET /items", GetItemsHandler(&items, &itemMU))
	mux.HandleFunc("POST /items/{id}/bids", CreateBidHandler(&items))
	mux.HandleFunc("GET /items/{id}/bids", GetBidsHandler(&items))
	mux.HandleFunc("GET /items/{id}/winner", GetItemWinner(&items))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server is running on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %s: %v\n", srv.Addr, err)
		}
	}()

	<-quit
	log.Println("Shutdown signal received, exiting...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting gracefully")
}
