package main

import (
	db "backend/auth/db"
	"backend/auth/handlers"
	mw "backend/middlewares"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	//getting environment variables
	port := os.Getenv("PORT")
	dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	//logger to inject to use across services
	logger := log.New(os.Stdout, "auth-service ", log.LstdFlags)
	bindAddress := fmt.Sprintf(":%s", port)

	//database init
	database, err := db.Initialize(dbUser, dbPassword, dbName)
	if err != nil {
		logger.Fatalf("Could not set up database: %v", err)
	}
	//close database connection on service end
	defer database.Conn.Close()

	//create handlers
	authHandler := handlers.NewAuth(logger, &database)

	//serve mux
	serveMux := mux.NewRouter()
	serveMux.Use(mw.JSONContentHeaders)
	serveMux.HandleFunc("/register", authHandler.Register).Methods("POST")
	serveMux.HandleFunc("/login", authHandler.Login).Methods("POST")
	serveMux.HandleFunc("/test", authHandler.Test).Methods("GET")

	// create a new server
	server := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      serveMux,          // set the default handler
		ErrorLog:     logger,            // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		logger.Printf("Starting server on port %s", port)

		err := server.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt)
	signal.Notify(stopChannel, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-stopChannel
	log.Println("Got signal:", sig)
	log.Println("Shutting down gracefully.")

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
	cancelFunc()
}
