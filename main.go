package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/cfg"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/clamav"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	config := &cfg.Config{}
	flag.StringVar(&config.ClamAVAddress, "clamav-address", "localhost", "ClamAV address to use")
	flag.IntVar(&config.ClamAVPort, "clamav-port", 3310, "ClamAV port to use")
	flag.Parse()

	fmt.Println("Server is starting...")

	router := http.NewServeMux()

	prometheus.MustRegister(clamav.NewCollector(*config))
	router.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", 8080),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	fmt.Println("Server is ready to handle requests at :", 8080)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %d: %v\n", 8080, err)
	}

	<-done
	fmt.Println("Server stopped")
}
