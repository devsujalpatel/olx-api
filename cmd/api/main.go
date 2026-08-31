package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok" }`))
	})

	err := http.ListenAndServe(":8090", mux)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
