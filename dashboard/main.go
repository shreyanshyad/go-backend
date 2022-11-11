package main

import (
	db "backend/dashboard/db"
	"backend/dashboard/handlers"
	"backend/dashboard/repository"
	"backend/dashboard/services"
	mw "backend/middlewares"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
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
	logger := log.New(os.Stdout, "dash-service ", log.LstdFlags)
	bindAddress := fmt.Sprintf(":%s", port)

	//database init
	database, err := db.Initialize(dbUser, dbPassword, dbName)
	if err != nil {
		logger.Fatalf("Could not set up database: %v", err)
	}
	//close database connection on service end
	defer database.Conn.Close()

	//redis init
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	//ping client
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf("Could not set up redis: %v", err)
	}
	logger.Println("Connected to redis")

	//create handlers
	dashRepo := repository.NewDashRepository(&database, logger)
	viewRepo := repository.NewViewRepository(&database, logger)
	roleRepo := repository.NewRoleRepository(&database, logger)
	roleService := services.NewRoleService(roleRepo, rdb, logger)
	viewService := services.NewViewService(viewRepo, roleService, logger)
	dashService := services.NewDashService(dashRepo, viewService, roleService, rdb, logger)

	dashHandler := handlers.NewDash(logger, dashService)

	//serve mux
	serveMux := mux.NewRouter()
	serveMux.Use(mw.LoggerMiddleware)
	serveMux.Use(mw.JSONContentHeaders) //adding content type to all responses
	serveMux.Use(mw.AuthMiddleware)

	//subrouter for post requests
	postR := serveMux.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/dashboard", dashHandler.CreateDash)
	postR.HandleFunc("/view", dashHandler.CreateView)
	postR.HandleFunc("/dashboard/{id}/users", dashHandler.AddUserToDash)
	postR.HandleFunc("/view/{id}/users", dashHandler.AddUserToView)

	//subrouter for delete requests
	deleteR := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/dashboard/{id}/users", dashHandler.DeleteUserFromDash)
	deleteR.HandleFunc("/view/{id}/users", dashHandler.DeleteUserFromView)

	//subrouter for patch requests
	patchR := serveMux.Methods(http.MethodPut).Subrouter()
	patchR.HandleFunc("/dashboard/{id}", dashHandler.UpdateDash)
	patchR.HandleFunc("/view/{id}", dashHandler.UpdateView)

	//subrouter for get requests
	getR := serveMux.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/dashboard", dashHandler.GetDashs)
	getR.HandleFunc("/dashboard/{id}", dashHandler.GetDash)
	getR.HandleFunc("/dashboard/{id}/users", dashHandler.GetUsersFromDash)
	getR.HandleFunc("/view/{id}", dashHandler.GetView)
	getR.HandleFunc("/roles", dashHandler.GetAllRoles)

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
