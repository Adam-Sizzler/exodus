package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/cerberus/subscription-page/backend/internal/config"
	"github.com/cerberus/subscription-page/backend/internal/server"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[FATAL] %v", err)
	}

	application, err := server.New(cfg)
	if err != nil {
		log.Fatalf("[FATAL] %v", err)
	}

	address := ":" + cfg.AppPort
	log.Printf("[CONFIG] CUSTOM_SUB_PREFIX: %s", config.DisplayPrefix(cfg.CustomSubPrefix))
	log.Printf("[CONFIG] listening on %s", address)

	httpServer := &http.Server{
		Addr:              address,
		Handler:           application,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("[FATAL] listen failed: %v", err)
	}
}
