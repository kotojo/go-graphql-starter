package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/kotojo/life-manager/generated"
	"github.com/kotojo/life-manager/middleware"
	"github.com/kotojo/life-manager/models"
	"github.com/kotojo/life-manager/resolvers"
	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("POSTGRESQL_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Postgresql connected")

	router := chi.NewRouter()
	router.Use(middleware.Http())

	usersService := models.NewUserService(db)
	router.Use(middleware.User(usersService))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{
		DB:               db,
		UsersService:     usersService,
		DocumentsService: models.NewDocumentService(db),
	}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
