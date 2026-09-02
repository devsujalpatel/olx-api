package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/devsujalpatel/olx-api/internal/config"
	"github.com/devsujalpatel/olx-api/internal/db"
	"github.com/devsujalpatel/olx-api/internal/handlers"
)

func main() {
	// checking if env is load
	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

  fmt.Println("database connected")
  fmt.Println("starting olx server...")

  // Creating mux or route
	mux := http.NewServeMux()

	// creating health check route
	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", handlers.List(db))

  // starting the http server
	 srv := http.Server{
		Addr: ":" + cfg.Port,
		Handler: mux,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout: time.Second * 60,
	}
	log.Printf("server is listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
