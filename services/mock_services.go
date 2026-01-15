package services

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/torresposso/gosmic/pb"
)

// MockPostService is a mock implementation of PostService
type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) List(ctx context.Context, client *pb.Client, query string) ([]pb.Post, error) {
	args := m.Called(ctx, client, query)
	posts, _ := args.Get(0).([]pb.Post)
	return posts, args.Error(1)
}

func (m *MockPostService) Get(ctx context.Context, client *pb.Client, id string) (*pb.Post, error) {
	args := m.Called(ctx, client, id)
	post, _ := args.Get(0).(*pb.Post)
	return post, args.Error(1)
}

func (m *MockPostService) Create(ctx context.Context, client *pb.Client, title, content string, isPublic bool) error {
	args := m.Called(ctx, client, title, content, isPublic)
	return args.Error(0)
}

func (m *MockPostService) Update(ctx context.Context, client *pb.Client, id string, title, content string, isPublic bool) error {
	args := m.Called(ctx, client, id, title, content, isPublic)
	return args.Error(0)
}

func (m *MockPostService) Delete(ctx context.Context, client *pb.Client, id string) error {
	args := m.Called(ctx, client, id)
	return args.Error(0)
}

func (m *MockPostService) TogglePublic(ctx context.Context, client *pb.Client, id string) error {
	args := m.Called(ctx, client, id)
	return args.Error(0)
}
