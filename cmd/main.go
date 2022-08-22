package main

import (
	"log"
	"net/http"
	"time"

	"github.com/itemun/crud-app/internal/repository/psql"
	"github.com/itemun/crud-app/internal/service"
	"github.com/itemun/crud-app/internal/transport/rest"
	"github.com/itemun/crud-app/pkg/database"

	_ "github.com/lib/pq"
)

func main() {
	// init db
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     "localhost",
		Port:     5433,
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "goLANGn1nja",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// init deps
	carsRepo := psql.NewCars(db)
	carsService := service.NewCars(carsRepo)
	handler := rest.NewHandler(carsService)

	// init & run server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
