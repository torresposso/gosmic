package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/torresposso/gosmic/pb"
	"github.com/torresposso/gosmic/repositories"
)

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(repositories.MockAuthRepository)
	service := NewAuthService(mockRepo)
	ctx := context.Background()
	client := &pb.Client{}

	t.Run("Success", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"
		expectedToken := "valid-token"

		mockRepo.On("Authenticate", ctx, client, email, password).Return(expectedToken, &pb.User{}, nil).Once()

		token, err := service.Login(ctx, client, email, password)

		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("MissingFields", func(t *testing.T) {
		token, err := service.Login(ctx, client, "", "password")
		assert.Error(t, err)
		assert.Equal(t, "email and password are required", err.Error())
		assert.Empty(t, token)
	})

	t.Run("RepoError", func(t *testing.T) {
		email := "test@example.com"
		password := "wrong-password"

		mockRepo.On("Authenticate", ctx, client, email, password).Return("", nil, errors.New("invalid credentials")).Once()

		token, err := service.Login(ctx, client, email, password)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authentication failed")
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(repositories.MockAuthRepository)
	service := NewAuthService(mockRepo)
	ctx := context.Background()
	client := &pb.Client{}

	t.Run("Success", func(t *testing.T) {
		email := "new@example.com"
		password := "securepassword123"
		name := "New User"

		mockRepo.On("CreateUser", ctx, client, mock.MatchedBy(func(data map[string]any) bool {
			return data["email"] == email && data["password"] == password && data["name"] == name
		})).Return(nil).Once()

		err := service.Register(ctx, client, email, password, name)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("MissingFields", func(t *testing.T) {
		err := service.Register(ctx, client, "", "password", "name")
		assert.Error(t, err)
		assert.Equal(t, "all fields are required", err.Error())
	})

	t.Run("ShortPassword", func(t *testing.T) {
		err := service.Register(ctx, client, "test@example.com", "short", "Name")
		assert.Error(t, err)
		assert.Equal(t, "password must be at least 8 characters", err.Error())
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.On("CreateUser", ctx, client, mock.Anything).Return(errors.New("db error")).Once()

		err := service.Register(ctx, client, "test@example.com", "password123", "Name")

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
