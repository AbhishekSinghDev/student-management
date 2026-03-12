package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AbhishekSinghDev/student-management/internal/config"
)

func startServer(server *http.Server) {
	slog.Info("server listening on", slog.String("address", server.Addr) )
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}

func main() {
	// load config
	config := config.MustLoad()

	// db setup
	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to go api"))
	})

	// setup server

	// buffered channel because we only need one signal to stop the server
	doneChannel := make(chan os.Signal, 1)

	signal.Notify(doneChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	server := http.Server{
		Addr: config.Address,
		Handler: router,
	}

	go startServer(&server)

	// passing any signal value in the done channel. this will hold the main function until all go routines stop completely
	<-doneChannel

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown the server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
