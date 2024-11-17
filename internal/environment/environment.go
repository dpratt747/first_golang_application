package environment

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetEnvVar(path string) (int, string, string, string, string, string) {

	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(path)
		if err != nil {
			message := fmt.Sprintf("Error loading .env file from path %v", path)
			log.Println(message)
		}
	}

	applicationPortString := getEnvOrDefault("APP_PORT", "8080")
	appPort, err := strconv.Atoi(applicationPortString)
	if err != nil {
		log.Println("Unable to convert appPortString to Int")
	}

	dbHost := getEnvOrDefault("DB_HOST", "localhost")

	runningMode := getEnvOrDefault("RUNNING_MODE", "locally")

	var dbPort string
	if runningMode == "docker" {
		dbPort = getEnvOrDefault("INTERNAL_DB_PORT", "5432")
	} else {
		dbPort = getEnvOrDefault("EXTERNAL_DB_PORT", "5432")
	}

	postgresUser := getEnvOrDefault("POSTGRES_USER", "postgres")

	postgresPassword := getEnvOrDefault("POSTGRES_PASSWORD", "postgres")

	postgresDb := getEnvOrDefault("POSTGRES_DB", "golang_db")

	return appPort, dbHost, dbPort, postgresUser, postgresPassword, postgresDb
}
