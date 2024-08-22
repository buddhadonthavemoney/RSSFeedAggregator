package main

import (
	"database/sql"
	"dbconnection/internal/database"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func init(){
	godotenv.Load()
}
func main(){
	portString := os.Getenv("PORT")
	dbUrl := os.Getenv("CONNECTION_STRING")
	if dbUrl == "" {
		log.Fatalf("CONNECTION_STRING variable not found")
	}

	conn, err:= sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	dbQueries := database.New(conn)
	apiCfg := apiConfig{
		DB: dbQueries,
	}


	if portString == ""{
		log.Fatalf("Port variable not found")
	}

	router := chi.NewRouter()

	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
			ExposedHeaders: []string{"Link"},
			AllowCredentials: false,
			MaxAge: 300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handleErr)
	v1Router.Get("/users", apiCfg.handlerGetUser)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":"+portString,
	}
	log.Printf("server starting on port %v", portString)

	log.Fatal(srv.ListenAndServe())
	if err != nil {
		log.Fatal(err)
	}
}
