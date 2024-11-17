package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	db "db_access/internal/database"
	"db_access/internal/domain"
	"db_access/internal/environment"
	"math/rand"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/stretchr/testify/assert"

	"github.com/pressly/goose/v3"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	containerPort string
	containerHost  string
	envPath string = "../../.env"
)

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func mustStartPostgresContainer() (func(context.Context) error, error) {

	_, _, _, postgresUser, postgresPassword, postgresDb := environment.GetEnvVar(envPath) 

	var (
		dbName = postgresDb
		dbPwd  = postgresPassword
		dbUser = postgresUser
	)

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		return dbContainer.Terminate, err
	}

	containerHost = dbHost
	containerPort = strings.ReplaceAll(string(dbPort), "tcp", "")
	containerPort = strings.ReplaceAll(containerPort, "/", "")

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not tear down postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	_, _, _, postgresUser, postgresPassword, postgresDb := environment.GetEnvVar(envPath) 

	dataSourceName := func(user, password, dbName, port, host string) string {
		return fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", user, password, dbName, port, host)
	}(postgresUser, postgresPassword, postgresDb, containerPort, containerHost)	

	srv := db.New(dataSourceName)
	assert.NotEqual(t, nil, srv, "New() returned nil")
}

func TestInsertNewUserSuccess(t *testing.T) {
	userForInsertion := domain.User{
		Username: "test user",
		Email:    "test@email.com",
	}

	_, _, _, postgresUser, postgresPassword, postgresDb := environment.GetEnvVar(envPath) 

	dataSourceName := func(user, password, dbName, port, host string) string {
		return fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", user, password, dbName, port, host)
	}(postgresUser, postgresPassword, postgresDb, containerPort, containerHost)	

	underTest := db.New(dataSourceName)

	sqlDb, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(sqlDb, "../../migrations"); err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		t.Log("Cleaning up after test")
		err := goose.DownTo(sqlDb, "../../migrations", 0)
		if err != nil {
			message := fmt.Sprintf("Error whilst cleaning migration: %v", err)
			t.Log(message)
		}
		sqlDb.Close()
	})

	userId, err := underTest.InsertNewUser(userForInsertion)
	assert.Equal(t, nil, err, "Some error occurred inserting the user. expected nil")

	var count int
	err = sqlDb.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 1, userId, fmt.Sprintf("Expected userId to equal to 1 got %v", userId))
	assert.Equal(t, 1, count, fmt.Sprintf("Expected count to equal to 1 got %v", count))
}

func TestInsertNewUserDuplicateUserEmailFailure(t *testing.T) {
	email := fmt.Sprintf("%v@email.com", randomString(10))

	userForInsertion1 := domain.User{
		Username: "test user1",
		Email:    email,
	}
	userForInsertion2 := domain.User{
		Username: "test user2",
		Email:    email,
	}

	_, _, _, postgresUser, postgresPassword, postgresDb := environment.GetEnvVar(envPath) 

	dataSourceName := func(user, password, dbName, port, host string) string {
		return fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", user, password, dbName, port, host)
	}(postgresUser, postgresPassword, postgresDb, containerPort, containerHost)		

	underTest := db.New(dataSourceName)

	sqlDb, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(sqlDb, "../../migrations"); err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		t.Log("Cleaning up after test")
		err := goose.DownTo(sqlDb, "../../migrations", 0)
		if err != nil {
			message := fmt.Sprintf("Error whilst cleaning migration: %v", err)
			t.Log(message)
		}
		sqlDb.Close()
	})

	_, err = underTest.InsertNewUser(userForInsertion1)
	assert.Equal(t, nil, err, "Some error occurred inserting the user. expected nil")
	_, err = underTest.InsertNewUser(userForInsertion2)
	_, isUniqueConstraintError := err.(*domain.UniqueConstraintDatabaseError)
	assert.True(t, isUniqueConstraintError, "Expected an UniqueConstraintDatabaseError when inserting a user with an already existing email address")
}

func TestGetAllUsersSuccess(t *testing.T) {
	userForInsertion1 := domain.User{
		Username: "test user 1",
		Email:    "email1@email.com",
	}

	userForInsertion2 := domain.User{
		Username: "test user 2",
		Email:    "email2@email.com",
	}

	_, _, _, postgresUser, postgresPassword, postgresDb := environment.GetEnvVar(envPath) 

	dataSourceName := func(user, password, dbName, port, host string) string {
		return fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", user, password, dbName, port, host)
	}(postgresUser, postgresPassword, postgresDb, containerPort, containerHost)	

	underTest := db.New(dataSourceName)

	sqlDb, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(sqlDb, "../../migrations"); err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		t.Log("Cleaning up after test")
		err := goose.DownTo(sqlDb, "../../migrations", 0)
		if err != nil {
			message := fmt.Sprintf("Error whilst cleaning migration: %v", err)
			t.Log(message)
		}
		sqlDb.Close()
	})

	_, err = underTest.InsertNewUser(userForInsertion1)
	assert.Equal(t, nil, err, "Some error occurred inserting the user. expected nil")
	_, err = underTest.InsertNewUser(userForInsertion2)
	assert.Equal(t, nil, err, "Some error occurred inserting the user. expected nil")

	getAllUsersResponse, _ := underTest.GetAllUsers()

	assert.Equal(t, 2, len(getAllUsersResponse), "expected GetAllUsers() to return a list of length equal to 2")
}

func TestGetAllUsersTombstoneSuccess(t *testing.T) {
	userForInsertion1 := domain.User{
		Username: "test user 1",
		Email:    "email1@email.com",
	}

	userForInsertion2 := domain.User{
		Username: "test user 2",
		Email:    "email2@email.com",
	}

	_, _, _, postgresUser, postgresPassword, postgresDb := environment.GetEnvVar(envPath) 

	dataSourceName := func(user, password, dbName, port, host string) string {
		return fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", user, password, dbName, port, host)
	}(postgresUser, postgresPassword, postgresDb, containerPort, containerHost)	
	underTest := db.New(dataSourceName)

	sqlDb, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(sqlDb, "../../migrations"); err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		t.Log("Cleaning up after test")
		err := goose.DownTo(sqlDb, "../../migrations", 0)
		if err != nil {
			message := fmt.Sprintf("Error whilst cleaning migration: %v", err)
			t.Log(message)
		}
		sqlDb.Close()
	})

	_, err = underTest.InsertNewUser(userForInsertion1)
	assert.Equal(t, nil, err, "Some error occurred inserting the user. expected nil")
	userId, err := underTest.InsertNewUser(userForInsertion2)
	assert.Equal(t, nil, err, "Some error occurred inserting the user. expected nil")

	// insert into tombstone with the userId above
	stmt := "INSERT INTO user_deletes(user_id) VALUES($1)"
	_, err = sqlDb.Exec(stmt, userId)
	if err != nil {
		log.Fatal(err)
	}
	getAllUsersResponse, _ := underTest.GetAllUsers()

	assert.Equal(t, 1, len(getAllUsersResponse), "expected GetAllUsers() to return a list of length equal to 1")
}

func TestSoftDeleteUserSuccess(t *testing.T) {
	userForInsertion := domain.User{
		Username: "test user 1",
		Email:    "email1@email.com",
	}

	_, _, _, postgresUser, postgresPassword, postgresDb := environment.GetEnvVar(envPath) 

	dataSourceName := func(user, password, dbName, port, host string) string {
		return fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", user, password, dbName, port, host)
	}(postgresUser, postgresPassword, postgresDb, containerPort, containerHost)	

	underTest := db.New(dataSourceName)

	sqlDb, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(sqlDb, "../../migrations"); err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		t.Log("Cleaning up after test")
		err := goose.DownTo(sqlDb, "../../migrations", 0)
		if err != nil {
			message := fmt.Sprintf("Error whilst cleaning migration: %v", err)
			t.Log(message)
		}
		sqlDb.Close()
	})

	userId, err := underTest.InsertNewUser(userForInsertion)
	assert.Equal(t, nil, err, "Some error occurred inserting the user. expected nil")
	err = underTest.SoftDeleteUser(userId)
	assert.Equal(t, nil, err, "Some error occurred inserting the user. expected nil")

	query := "SELECT COUNT(*) FROM user_deletes ud WHERE ud.user_id = $1"
	var count int
	err = sqlDb.QueryRow(query, userId).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 1, count, "expected SoftDeleteUser() to persist 1 row to the user_deletes table")
}
