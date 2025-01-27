/*
Copyright 2020 Christian Niehoff.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/clamav"
	"github.com/r3kzi/clamav-prometheus-exporter/pkg/collector"
	log "github.com/sirupsen/logrus"
)

var (
	address string
	port    int
	network string
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	flag.StringVar(&address, "clamav-address", "localhost", "ClamAV address to use")
	flag.IntVar(&port, "clamav-port", 3310, "ClamAV port to use")
	flag.StringVar(&network, "network", "tcp", "Network mode to use, typically tcp or unix (socket)")
	flag.Parse()
}

func main() {
	log.Info("Server is starting...")

	if strings.EqualFold(network, "tcp") {
		address = fmt.Sprintf("%s:%d", address, port)
	}

	client := clamav.New(address, network)
	coll := collector.New(*client)
	prometheus.MustRegister(coll)

	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", 9810),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)

	quit := make(chan os.Signal, 1)
	// catch SIGTERM or SIGINT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	log.Info("Server is ready to handle requests at :", 9810)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %d: %v\n", 9810, err)
	}

	<-done
	log.Info("Server stopped")
}
