package repositories

import (
	"context"

	"github.com/torresposso/gosmic/pb"
)

// PostRepository defines the interface for post data access
type PostRepository interface {
	List(ctx context.Context, client *pb.Client) ([]pb.Post, error)
	Get(ctx context.Context, client *pb.Client, id string) (*pb.Post, error)
	Create(ctx context.Context, client *pb.Client, title, content string, isPublic bool) error
	Update(ctx context.Context, client *pb.Client, id string, data map[string]any) error
	Delete(ctx context.Context, client *pb.Client, id string) error
	TogglePublic(ctx context.Context, client *pb.Client, id string) error
}

// PBPostRepository implements PostRepository using PocketBase
type PBPostRepository struct{}

func NewPostRepository() PostRepository {
	return &PBPostRepository{}
}

func (r *PBPostRepository) List(ctx context.Context, client *pb.Client) ([]pb.Post, error) {
	return client.ListPosts()
}

func (r *PBPostRepository) Get(ctx context.Context, client *pb.Client, id string) (*pb.Post, error) {
	return client.GetPost(id)
}

func (r *PBPostRepository) Create(ctx context.Context, client *pb.Client, title, content string, isPublic bool) error {
	return client.CreatePost(title, content, isPublic)
}

func (r *PBPostRepository) Update(ctx context.Context, client *pb.Client, id string, data map[string]any) error {
	return client.UpdatePost(id, data)
}

func (r *PBPostRepository) Delete(ctx context.Context, client *pb.Client, id string) error {
	return client.DeletePost(id)
}

func (r *PBPostRepository) TogglePublic(ctx context.Context, client *pb.Client, id string) error {
	post, err := client.GetPost(id)
	if err != nil {
		return err
	}
	return client.UpdatePost(id, map[string]any{
		"public": !post.Public,
	})
}
