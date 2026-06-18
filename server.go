package main

import (
	"context"
	"fitquest-backend/graph"
	"log"
	"net/http"
	"os"
	"time"

	"fitquest-backend/auth"
	"fitquest-backend/database"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	// Load the .env file from the current directory
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: No .env file found, relying on system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbClient := database.ConnectDB()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := dbClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			DB: dbClient,
		},
	}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", auth.Middleware()(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
