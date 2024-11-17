package environment

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetEnvVar(path string) (int, string, string, string, string, string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Println("Error loading .env file")
	}

	applicationPortString := getEnvOrDefault("APP_PORT", "8080")
	appPort, err := strconv.Atoi(applicationPortString)
	if err != nil {
		log.Println("Unable to convert appPortString to Int")
	}

	dbHost := getEnvOrDefault("DB_HOST", "localhost")

	dbPort := getEnvOrDefault("DB_PORT", "5432")

	postgresUser := getEnvOrDefault("POSTGRES_USER", "postgres")

	postgresPassword := getEnvOrDefault("POSTGRES_PASSWORD", "postgres")

	postgresDb := getEnvOrDefault("POSTGRES_DB", "golang_db")

	return appPort, dbHost, dbPort, postgresUser, postgresPassword, postgresDb
}
