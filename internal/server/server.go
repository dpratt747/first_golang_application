package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"db_access/internal/database"
	"db_access/internal/enums"
)

type Server struct {
	Port int
	Db database.DatabaseService
}

func New() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))


	dataSourceName := func(user, password, dbName, port, host string) string {
		return fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", user, password, dbName, port, host)
	}(string(enums.Username), string(enums.Password), string(enums.Database), string(enums.Port), string(enums.Host))

	NewServer := &Server{
		Port: port,
		Db: database.New(dataSourceName),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
