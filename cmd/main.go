package main

import (
	"github.com/itemun/crud-app/internal/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/itemun/crud-app/internal/repository/psql"
	"github.com/itemun/crud-app/internal/service"
	"github.com/itemun/crud-app/internal/transport/rest"
	"github.com/itemun/crud-app/pkg/database"

	_ "github.com/lib/pq"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	cfg, err := config.New("configs", "example")
	if err != nil {
		logrus.Fatal(err)
	}

	// init db
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Username: cfg.DBUser,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
		Password: cfg.DBPassword,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	// init deps
	carsRepo := psql.NewCars(db)
	carsService := service.NewCars(carsRepo)
	handler := rest.NewHandler(carsService)

	// init & run server
	srv := &http.Server{
		Addr:    ":" + cfg.SrvPort,
		Handler: handler.InitRouter(),
	}

	logrus.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatal(err)
	}
}
