package server

import (
	"fmt"
	"log"
	"net/http"

	// "os"
	// "strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"db_access/internal/database"
	"db_access/internal/enums"
)

type Server struct {
	Port int
	Db   database.DatabaseService
}

func New() *http.Server {
	// todo: make these either config values or environment variables
	port := 8080
	host := "127.0.0.1"

	dataSourceName := func(user, password, dbName, port, host string) string {
		return fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", user, password, dbName, port, host)
	}(string(enums.Username), string(enums.Password), string(enums.Database), string(enums.Port), string(enums.Host))

	db := database.New(dataSourceName)

	NewServer := &Server{
		Port: port,
		Db:   db,
	}

	address := fmt.Sprintf("%v:%d", host, NewServer.Port)

	message := fmt.Sprintf("Server has started on: %v", address)
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
