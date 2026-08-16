package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fullcyclelabs.com/weather-app/internal/clients/viacep"
	"fullcyclelabs.com/weather-app/internal/clients/weatherapi"
	"fullcyclelabs.com/weather-app/internal/config"
	"fullcyclelabs.com/weather-app/internal/handler"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (local development)
	// In Cloud Run, env vars are passed directly
	_ = godotenv.Load(".env")

	cfg := config.Load()

	viacepClient := viacep.NewClient(http.DefaultClient)
	weatherapiClient := weatherapi.NewClient(http.DefaultClient, cfg.WeatherApiKey)
	weatherHandler := handler.NewWeatherHandler(viacepClient, weatherapiClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", weatherHandler.HandleWeatherRequest)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: mux,
	}

	go func() {
		log.Printf("Server starting on %s...", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("\nShutdown signal received. Cleaning up...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly.")
}
