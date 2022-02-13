package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	hub := newHub()
	go hub.run(ctx, &wg)
	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	}).Methods(http.MethodGet)
	srv := http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		<-signals
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
		cancel()
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	wg.Wait()
}
