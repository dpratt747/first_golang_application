package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"db_access/internal/domain"

	sv "db_access/internal/server"
	testMocks "db_access/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	expectedStatusCode := http.StatusOK
	assert.Equal(t, expectedStatusCode, rr.Code, fmt.Sprintf("Expected response status to equal %v. [actual]: %v", expectedStatusCode, rr.Code))
	expected := "{\"message\":\"Hello World\"}"
	assert.Equal(t, expected, rr.Body.String(), fmt.Sprintf("Expected response body to equal %v. [actual]: %v", expected, rr.Body.String()))
}

func TestGetAllUsersSuccess(t *testing.T) {
	user := domain.User{
		ID: 0,
		Username: "New User",
		Email: "NewEmail@github.com",
	}

	service := new(testMocks.MockDBService)

	userList := []domain.User {user}

	service.On("GetAllUsers").Return(userList, nil)

	s := &sv.Server{
		Port: 8080,
		Db: service,
	}
	r := gin.New()
	r.GET("/users", s.GetAllUsersHandler)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	expectedStatusCode := http.StatusOK
	assert.Equal(t, expectedStatusCode, rr.Code, fmt.Sprintf("Expected response status to equal %v. [actual]: %v", expectedStatusCode, rr.Code))
	expected := `[{"id":0,"username":"New User","email":"NewEmail@github.com"}]`
	assert.Equal(t, expected, rr.Body.String(), fmt.Sprintf("Expected response body to equal %v. [actual]: %v", expected, rr.Body.String()))
}

func TestGetAllUsersFailure(t *testing.T) {
	service := new(testMocks.MockDBService)

	service.On("GetAllUsers").Return([]domain.User{}, errors.New("Something went wrong"))

	s := &sv.Server{
		Port: 8080,
		Db: service,
	}
	r := gin.New()
	r.GET("/users", s.GetAllUsersHandler)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	expectedStatusCode := http.StatusInternalServerError
	assert.Equal(t, expectedStatusCode, rr.Code, fmt.Sprintf("Expected response status to equal %v. [actual]: %v", expectedStatusCode, rr.Code))
	expected := "{}"
	assert.Equal(t, expected, rr.Body.String(), fmt.Sprintf("Expected response body to equal %v. [actual]: %v", expected, rr.Body.String()))
}

func TestDeleteUserHandlerSuccess(t *testing.T) {
	userDeletion := domain.UserDeletion{
		UserId: 1,
	}

	service := new(testMocks.MockDBService)
	service.On("SoftDeleteUser", mock.Anything).Return(nil)

	s := &sv.Server{
		Port: 8080,
		Db: service,
	}
	r := gin.New()
	r.DELETE("/user", s.DeleteUserHandler)

	jsonData, err := json.Marshal(userDeletion)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}

	// Create a test HTTP request
	req, err := http.NewRequest("DELETE", "/user", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	expectedStatusCode := http.StatusNoContent
	assert.Equal(t, expectedStatusCode, rr.Code, fmt.Sprintf("Expected response status to equal %v. [actual]: %v", expectedStatusCode, rr.Code))
	assert.Empty(t, rr.Body.String(), fmt.Sprintf("Expected response body to be empty. [actual]: %v", rr.Body.String()))
}

func TestDeleteUserHandlerUniqueConstraintFailure(t *testing.T) {
	userDeletion := domain.UserDeletion{
		UserId: 1,
	}

	service := new(testMocks.MockDBService)
	service.On("SoftDeleteUser", mock.Anything).Return(&domain.UniqueConstraintDatabaseError{Message: "some issue"})

	s := &sv.Server{
		Port: 8080,
		Db: service,
	}
	r := gin.New()
	r.DELETE("/user", s.DeleteUserHandler)

	jsonData, err := json.Marshal(userDeletion)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}

	// Create a test HTTP request
	req, err := http.NewRequest("DELETE", "/user", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	expectedStatusCode := http.StatusBadRequest
	assert.Equal(t, expectedStatusCode, rr.Code, fmt.Sprintf("Expected response status to equal %v. [actual]: %v", expectedStatusCode, rr.Code))
	expected := "{\"error\":\"Unable to delete this user as they have already been deleted\"}"
	assert.Equal(t, expected, rr.Body.String(), fmt.Sprintf("Expected response body to equal %v. [actual]: %v", expected, rr.Body.String()))
}

func TestDeleteUserHandlerUserNotFoundFailure(t *testing.T) {
	userDeletion := domain.UserDeletion{
		UserId: 1,
	}

	service := new(testMocks.MockDBService)
	service.On("SoftDeleteUser", mock.Anything).Return(&domain.UserNotFoundError{Message: "some issue"})

	s := &sv.Server{
		Port: 8080,
		Db: service,
	}
	r := gin.New()
	r.DELETE("/user", s.DeleteUserHandler)

	jsonData, err := json.Marshal(userDeletion)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}

	// Create a test HTTP request
	req, err := http.NewRequest("DELETE", "/user", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	expectedStatusCode := http.StatusBadRequest
	assert.Equal(t, expectedStatusCode, rr.Code, fmt.Sprintf("Expected response status to equal %v. [actual]: %v", expectedStatusCode, rr.Code))
	expected := "{\"error\":\"Unable to delete this user as they do not exist\"}"
	assert.Equal(t, expected, rr.Body.String(), fmt.Sprintf("Expected response body to equal %v. [actual]: %v", expected, rr.Body.String()))
}

func TestInsertNewUserHandlerSuccess(t *testing.T) {
	user := domain.User{
		ID: 0,
		Username: "New User",
		Email: "NewEmail@github.com",
	}

	service := new(testMocks.MockDBService)

	service.On("InsertNewUser", user).Return(10, nil)

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

	expectedStatusCode := http.StatusCreated
	assert.Equal(t, expectedStatusCode, rr.Code, fmt.Sprintf("Expected response status to equal %v. [actual]: %v", expectedStatusCode, rr.Code))
	expected := "{\"userId\":10}"
	assert.Equal(t, expected, rr.Body.String(), fmt.Sprintf("Expected response body to equal %v. [actual]: %v", expected, rr.Body.String()))
}

func TestInsertNewUserHandlerDuplicateEmailAddressFailure(t *testing.T) {
	user := domain.User{
		ID: 0,
		Username: "New User",
		Email: "NewEmail@github.com",
	}

	service := new(testMocks.MockDBService)

	uniqueRequestError := &domain.UniqueConstraintDatabaseError{Message: "This email is not unique"}

	service.On("InsertNewUser", user).Return(0, uniqueRequestError)

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

	expectedStatusCode := http.StatusBadRequest
	assert.Equal(t, expectedStatusCode, rr.Code, fmt.Sprintf("Expected response status to equal %v. [actual]: %v", expectedStatusCode, rr.Code))
	expected := "{\"error\":\"cannot insert user as this email is already used\"}"
	assert.Equal(t, expected, rr.Body.String(), fmt.Sprintf("Expected response body to equal %v. [actual]: %v", expected, rr.Body.String()))
}

func TestInsertNewUserHandlerFailureStatusCode422(t *testing.T) {
	service := new(testMocks.MockDBService)

	service.AssertNotCalled(t, "InsertNewUser", mock.Anything)
	
	s := &sv.Server{
		Port: 8080,
		Db: service,
	}
	r := gin.New()
	r.POST("/user", s.InsertNewUserHandler)

	invalidJsonString := `{ "invalid": "unprocessable" }`

	jsonData, err := json.Marshal(invalidJsonString)
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

	expectedStatusCode := http.StatusUnprocessableEntity
	assert.Equal(t, expectedStatusCode, rr.Code, fmt.Sprintf("Expected response status to equal %v. [actual]: %v", expectedStatusCode, rr.Code))
	expected := `{"error":"json: cannot unmarshal string into Go value of type domain.User"}`
	assert.Equal(t, expected, rr.Body.String(), fmt.Sprintf("Expected response body to equal %v. [actual]: %v", expected, rr.Body.String()))
}