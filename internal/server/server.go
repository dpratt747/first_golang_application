package server

import (
	"fmt"
	"log"
	"net/http"

	"time"

	_ "github.com/joho/godotenv/autoload"

	"db_access/internal/database"
	"db_access/internal/environment"
)

type Server struct {
	Port int
	Db   database.DatabaseService
}

func New() *http.Server {

	appPort, dbHost, dbPort, postgresUser, postgresPassword, postgresDb := environment.GetEnvVar(".env")

	dataSourceName := func(user, password, dbName, port, host string) string {
		return fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", user, password, dbName, dbPort, host)
	}(postgresUser, postgresPassword, postgresDb, dbPort, dbHost)

	db := database.New(dataSourceName)
	message := fmt.Sprintf("Database connection on: %v", dataSourceName)
	log.Println(message)

	NewServer := &Server{
		Port: appPort,
		Db:   db,
	}

	address := fmt.Sprintf(":%d", NewServer.Port)

	message = fmt.Sprintf("Server has started on: %v", address)
	log.Println(message)

	server := &http.Server{
		Addr:         address,
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
