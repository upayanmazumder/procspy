package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("procspy is live!"))
	})

	r.HandleFunc("/api/v1/machines/register", RegisterMachine).Methods("POST")
	r.HandleFunc("/api/v1/metrics", PushMetrics).Methods("POST")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("procspy server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}
	log.Println("server exited properly")
}

// placeholder handlers
func RegisterMachine(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"registered"}`))
}

func PushMetrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"metrics_received"}`))
}
