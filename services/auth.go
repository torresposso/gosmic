package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/torresposso/gosmic/pb"
	"github.com/torresposso/gosmic/repositories"
)

type AuthService interface {
	Login(ctx context.Context, client *pb.Client, email, password string) (string, error)
	Register(ctx context.Context, client *pb.Client, email, password, name string) error
}

type authService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(ctx context.Context, client *pb.Client, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email and password are required")
	}

	token, _, err := s.repo.Authenticate(ctx, client, email, password)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	return token, nil
}

func (s *authService) Register(ctx context.Context, client *pb.Client, email, password, name string) error {
	if email == "" || password == "" || name == "" {
		return errors.New("all fields are required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	data := map[string]any{
		"email":           email,
		"password":        password,
		"passwordConfirm": password, // Assuming confirmation is handled at the handler level or we just reuse password
		"name":            name,
	}

	return s.repo.CreateUser(ctx, client, data)
}
