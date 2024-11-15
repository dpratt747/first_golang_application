package tests

import (
	"db_access/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockDBService struct {
	mock.Mock
}

func (ms *MockDBService) Health() map[string]string {
	args := ms.Called()
	response := map[string]string {
		"key1": args.String(0),
	}
	return response
}
func (ms *MockDBService) InsertNewUser(user domain.User) (int, error) {
	args := ms.Called(user)
	return args.Int(0), args.Error(1)
}