package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/api/routes"
	"server/internal/models"
	"server/internal/security"
	"server/internal/store"
	"time"

	"github.com/gorilla/handlers"
)

func main() {

	var config = models.LoadConfig("config.yaml")

	mClient, err := store.Connect(config.URI())
	if err != nil {
		log.Printf("[db] failed to connect: %v", err)
		os.Exit(1)
	}

	defer func() {
		if err := mClient.Disconnect(context.Background()); err != nil {
			log.Printf("[db] error disconnecting: %v", err)
		}
	}()

	crypto, err := security.NewAESGCM(security.EncryptionKey())
	if err != nil {
		log.Printf("[error] failed to initialize aesgcm: %v", err)
		os.Exit(1)
	}

	db := mClient.Database(config.Database.Name)
	router, err := routes.SetupRoutes(db, crypto, config)
	if err != nil {
		log.Printf("[server] failed to setup routes: %v", err)
		os.Exit(1)
	}

	// Initialize the server
	srv := &http.Server{

		// CORS...
		Handler: handlers.CORS(
			handlers.AllowedOrigins(config.AllowedOrigins()),
			handlers.AllowedMethods([]string{
				"GET", "POST", "HEAD", "OPTIONS", "PUT",
			}),
			handlers.AllowedHeaders([]string{
				"X-Requested-With", "Content-Type", "Authorization",
			}),
			handlers.AllowCredentials(),
		)(router),
		Addr: config.Address(),

		// Enforce timeouts for server
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		<-stop

		log.Println("[server] closing server")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("[server] [err] closing server: %v", err)
		}
	}()

	log.Printf("[server] started at %s", config.Address())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[server] [err]: %v", err)
	}

	log.Println("[server] closed successfully")
}
