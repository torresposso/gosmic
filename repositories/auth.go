package repositories

import (
	"context"

	"github.com/torresposso/gosmic/pb"
)

// AuthRepository defines the interface for authentication data access
type AuthRepository interface {
	Authenticate(ctx context.Context, client *pb.Client, email, password string) (string, *pb.User, error)
	CreateUser(ctx context.Context, client *pb.Client, data map[string]any) error
}

// PBAuthRepository implements AuthRepository using PocketBase
type PBAuthRepository struct{}

func NewAuthRepository() AuthRepository {
	return &PBAuthRepository{}
}

func (r *PBAuthRepository) Authenticate(ctx context.Context, client *pb.Client, email, password string) (string, *pb.User, error) {
	return client.AuthWithPassword(email, password)
}

func (r *PBAuthRepository) CreateUser(ctx context.Context, client *pb.Client, data map[string]any) error {
	return client.CreateRecord("users", data)
}
