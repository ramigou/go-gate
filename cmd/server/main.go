package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ramigou/go-gate/internal/config"
	"github.com/ramigou/go-gate/internal/proxy"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	proxyServer := proxy.NewServer(cfg)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: proxyServer.Handler(),
	}

	go func() {
		log.Printf("Starting reverse proxy server on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}