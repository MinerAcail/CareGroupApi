package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/kobbi/vbciapi/database"
	"github.com/kobbi/vbciapi/graph"
	"github.com/kobbi/vbciapi/graph/model"
	"github.com/kobbi/vbciapi/jwt/middleware"
	"gorm.io/gorm"
)

const defaultPort = "9091"

func main() {
	config := &database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "justtest",
		Password: "0000",
		DBName:   "testdb1",
		SSLMode:  "disable",
	}

	db, err := database.NewConnection(config)
	if err != nil {
		// Handle the error
		log.Fatal(err)
	}

	err = PerformMigrations(db)
	if err != nil {
		// Handle the error
		log.Fatal(err)
	}
	err = MigrateInfor(db)
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
		AllowedOrigins:   []string{"http://localhost:5173"}, // Set to a slice of allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	router.Use(corsConfig.Handler) // Use the cors handler

	router.Use(middleware.Middleware)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// PerformMigrations performs database migrations using Gorm's AutoMigrate function.
func PerformMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&model.Member{}, &model.Church{}, &model.SubChurch{}, &model.CallCenter{}, &model.Registration{}, &model.MigrationRequest{}, &model.RegistrationByCallAgent{}, &model.JobInfo{}, &model.EmergencyContact{})
	if err != nil {
		return err
	}

	return nil
}

func MigrateInfor(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.FamilyInfo{}, &model.ChurchMinistryRole{}, &model.MemberChurchMinistryRole{}, &model.MemberChildren{}); err != nil {
		return err
	}
	return nil
}
