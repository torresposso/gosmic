package repositories

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/torresposso/gosmic/pb"
)

// MockAuthRepository is a mock implementation of AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Authenticate(ctx context.Context, client *pb.Client, email, password string) (string, *pb.User, error) {
	args := m.Called(ctx, client, email, password)
	user, _ := args.Get(1).(*pb.User)
	return args.String(0), user, args.Error(2)
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, client *pb.Client, data map[string]any) error {
	args := m.Called(ctx, client, data)
	return args.Error(0)
}

// MockPostRepository is a mock implementation of PostRepository
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) List(ctx context.Context, client *pb.Client) ([]pb.Post, error) {
	args := m.Called(ctx, client)
	posts, _ := args.Get(0).([]pb.Post)
	return posts, args.Error(1)
}

func (m *MockPostRepository) Get(ctx context.Context, client *pb.Client, id string) (*pb.Post, error) {
	args := m.Called(ctx, client, id)
	post, _ := args.Get(0).(*pb.Post)
	return post, args.Error(1)
}

func (m *MockPostRepository) Create(ctx context.Context, client *pb.Client, title, content string, isPublic bool) error {
	args := m.Called(ctx, client, title, content, isPublic)
	return args.Error(0)
}

func (m *MockPostRepository) Update(ctx context.Context, client *pb.Client, id string, data map[string]any) error {
	args := m.Called(ctx, client, id, data)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(ctx context.Context, client *pb.Client, id string) error {
	args := m.Called(ctx, client, id)
	return args.Error(0)
}

func (m *MockPostRepository) TogglePublic(ctx context.Context, client *pb.Client, id string) error {
	args := m.Called(ctx, client, id)
	return args.Error(0)
}
