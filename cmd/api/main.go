package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/devsujalpatel/olx-api/internal/config"
	"github.com/devsujalpatel/olx-api/internal/db"
	"github.com/devsujalpatel/olx-api/internal/handlers"
	"github.com/devsujalpatel/olx-api/internal/middleware"
)

func main() {
	// checking if env is load
	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	// Logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

  fmt.Println("database connected")
  fmt.Println("starting olx server...")

  lh := handlers.NewListingHandler(db, logger)

  // Creating mux or route
	mux := http.NewServeMux()

	// creating health check route
	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	handler := middleware.RequestId(mux)

  // starting the http server
	 srv := http.Server{
		Addr: ":" + cfg.Port,
		Handler: handler,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout: time.Second * 60,
	}
	
	log.Printf("server is listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
