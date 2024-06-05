package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"http/database"
	"http/handlers"
	"http/middlewares"
	"http/pb"
	"http/repository"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[Main Service] Error loading .env file")
	}
	authConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[Main Service] did not connect: %v", err)
	}
	defer authConn.Close()
	authClient := pb.NewAuthServiceClient(authConn)

	emailConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[Main Service] did not connect: %v", err)
	}
	defer emailConn.Close()
	emailClient := pb.NewEmailServiceClient(emailConn)
	connectionString := os.Getenv("DATABASE_URL")
	fmt.Println("DATABASE_URL: ", connectionString)
	db, err := database.NewDatabase(connectionString)
	if err != nil {
		log.Fatalf("[Main Service] error on connecting database: %v", err)

	}
	repo := repository.NewRepository(db.DB)
	s := handlers.NewHandler(authClient, emailClient, *repo)

	r := mux.NewRouter()
	r.HandleFunc("/auth/login", s.Login).Methods("POST")
	r.HandleFunc("/auth/register", s.Register).Methods("POST")
	r.HandleFunc("/books", s.GetBooks).Methods("GET")
	r.HandleFunc("/books", s.CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", s.GetBook).Methods("GET")
	r.HandleFunc("/books/{id}", s.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", s.DeleteBook).Methods("DELETE")
	r.HandleFunc("/verify", s.VerifyCode).Methods("POST")
	r.Use(middlewares.LoggingMiddleware)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("[Main Service] Starting server at :8080")
	log.Fatal(srv.ListenAndServe())
}
