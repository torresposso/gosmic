package services

import (
	"context"
	"strings"

	"github.com/torresposso/gosmic/pb"
	"github.com/torresposso/gosmic/repositories"
)

type PostService interface {
	List(ctx context.Context, client *pb.Client, query string) ([]pb.Post, error)
	Get(ctx context.Context, client *pb.Client, id string) (*pb.Post, error)
	Create(ctx context.Context, client *pb.Client, title, content string, isPublic bool) error
	Update(ctx context.Context, client *pb.Client, id string, title, content string, isPublic bool) error
	Delete(ctx context.Context, client *pb.Client, id string) error
	TogglePublic(ctx context.Context, client *pb.Client, id string) error
}

type postService struct {
	repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) List(ctx context.Context, client *pb.Client, query string) ([]pb.Post, error) {
	posts, err := s.repo.List(ctx, client)
	if err != nil {
		return nil, err
	}

	if query != "" {
		filtered := []pb.Post{}
		q := strings.ToLower(query)
		for _, p := range posts {
			if strings.Contains(strings.ToLower(p.Title), q) || strings.Contains(strings.ToLower(p.Content), q) {
				filtered = append(filtered, p)
			}
		}
		posts = filtered
	}

	return posts, nil
}

func (s *postService) Get(ctx context.Context, client *pb.Client, id string) (*pb.Post, error) {
	return s.repo.Get(ctx, client, id)
}

func (s *postService) Create(ctx context.Context, client *pb.Client, title, content string, isPublic bool) error {
	return s.repo.Create(ctx, client, title, content, isPublic)
}

func (s *postService) Update(ctx context.Context, client *pb.Client, id string, title, content string, isPublic bool) error {
	data := map[string]any{
		"title":   title,
		"content": content,
		"public":  isPublic,
	}
	return s.repo.Update(ctx, client, id, data)
}

func (s *postService) Delete(ctx context.Context, client *pb.Client, id string) error {
	return s.repo.Delete(ctx, client, id)
}

func (s *postService) TogglePublic(ctx context.Context, client *pb.Client, id string) error {
	return s.repo.TogglePublic(ctx, client, id)
}
