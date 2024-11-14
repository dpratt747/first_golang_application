package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"db_access/internal/domain"

	sv "db_access/internal/server"
	"db_access/tests/mocks"

	"github.com/gin-gonic/gin"
)

func TestHelloWorldHandler(t *testing.T) {
	s := &sv.Server{}
	r := gin.New()
	r.GET("/", s.HelloWorldHandler)
	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check the response body
	expected := "{\"message\":\"Hello World\"}"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestInsertNewUserHandler(t *testing.T) {

	user := domain.User{
		ID: 0,
		Username: "New User",
		Email: "NewEmail@github.com",
	}

	service := new(tests.MockDBService)

	service.On("InsertNewUser", user).Return(10)

	s := &sv.Server{
		Port: 8080,
		Db: service,
	}
	r := gin.New()
	r.POST("/user", s.InsertNewUserHandler)

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}

	// Create a test HTTP request
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check the response body
	expected := "10"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}