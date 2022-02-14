package main

import (
	"context"
	"embed"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	//go:embed index.html
	view embed.FS
)

func main() {
	quit := make(chan struct{})
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.FS(view)))
	r.HandleFunc("/", serveView).Methods(http.MethodGet)
	r.HandleFunc("/{channel}/ws", serveWS).Methods(http.MethodGet)
	srv := http.Server{
		Handler:      r,
		Addr:         ":80",
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
		close(quit)
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-quit
}

func serveView(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &Client{
		conn:    conn,
		channel: mux.Vars(r)["channel"],
	}
	go client.writePump()
	go client.readPump()
}
