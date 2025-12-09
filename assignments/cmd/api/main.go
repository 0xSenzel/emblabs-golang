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

	"github.com/0xsenzel/emblabs-golang/internal/handler"
	"github.com/0xsenzel/emblabs-golang/internal/service"
)

func main() {
	paymentService := service.NewPaymentService()

	// inject service dependency
	h := handler.NewHandler(paymentService)

	// register routes
	http.HandleFunc("/pay", h.Pay)

	// configure server with timeouts
	port := ":8080"
	server := &http.Server{
		Addr:         port,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// start server in background
	go func() {
		fmt.Printf("Starting Payment Service on http://localhost%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// wait for interrupt signal (CTRL+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// graceful shutdown with 30 second timeout
	fmt.Println("\nShutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Payment Service stopped gracefully")
}
