package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors" // Import the cors package
	"github.com/kobbi/vbciapi/database"
	"github.com/kobbi/vbciapi/graph"
	"github.com/kobbi/vbciapi/graph/model"
	"github.com/kobbi/vbciapi/jwt/middleware"
	"gorm.io/gorm"
)

const defaultPort = "9090"

func main() {
	config := &database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "justtest",
		Password: "0000",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	db, err := database.NewConnection(config)
	if err != nil {
		// Handle the error
		log.Fatal(err)
	}

	// Perform database migrations
	err = PerformMigrations(db)
	if err != nil {
		// Handle the error
		log.Fatal(err)
	}

	resolver := &graph.Resolver{
		DB: db,
	}

	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	router := chi.NewRouter()

	// Configure CORS
	corsConfig := cors.New(cors.Options{
		// { "http://localhost:5173", "http://yourdomain.com", "https://anotherdomain.com"},
		AllowedOrigins:   []string{"*"}, // You can configure specific origins here
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	router.Use(corsConfig.Handler) // Use the cors handler

	router.Use(middleware.AuthenticationMiddleware)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// PerformMigrations performs database migrations using Gorm's AutoMigrate function.
func PerformMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&model.Leader{}, &model.Student{}, &model.Registration{})
	if err != nil {
		return err
	}

	return nil
}
