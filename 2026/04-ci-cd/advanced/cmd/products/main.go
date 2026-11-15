package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"products/internal/infrastucture/handlers"
	"products/internal/infrastucture/metrics"
	"products/internal/infrastucture/postgres"
	"products/internal/interfaces/repository"
	"products/internal/usecases/storage"
)

func main() {
	var config = &postgres.Config{
		IP:       os.Getenv("POSTGRES_ADDR"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASS"),
		Database: os.Getenv("POSTGRES_DB"),
	}

	dbClient, err := postgres.New(config)
	if err != nil {
		log.Fatalln("Error during Postgres initialization")
	}

	productStorage := storage.New(repository.New(dbClient))

	router := mux.NewRouter()
	handlers.Make(router, productStorage)

	r := prometheus.NewRegistry()
	r.MustRegister(metrics.HttpRequestTotal)
	r.MustRegister(metrics.HttpRequestsDurationHistogram)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	router.Handle("/metrics", handler)

	srv := &http.Server{
		Addr:    ":30001",
		Handler: router,
	}

	go func() {
		listener := make(chan os.Signal, 1)
		signal.Notify(listener, os.Interrupt, syscall.SIGTERM)
		fmt.Println("Received a shutdown signal:", <-listener)

		if err := srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			fmt.Println("Failed to gracefully shutdown ", err)
		}
	}()

	log.Print("[*]  Listening...")

	if err := srv.ListenAndServe(); err != nil {
		log.Print("Failed to listen and serve ", err)
	}

	log.Print("Server shutdown")
}
