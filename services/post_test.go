package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/torresposso/gosmic/pb"
	"github.com/torresposso/gosmic/repositories"
)

func TestPostService_List(t *testing.T) {
	mockRepo := new(repositories.MockPostRepository)
	service := NewPostService(mockRepo)
	ctx := context.Background()
	client := &pb.Client{}

	posts := []pb.Post{
		{ID: "1", Title: "First Post", Content: "Hello world", Public: true},
		{ID: "2", Title: "Secret Post", Content: "Classified info", Public: false},
	}

	t.Run("SuccessNoQuery", func(t *testing.T) {
		mockRepo.On("List", ctx, client).Return(posts, nil).Once()

		result, err := service.List(ctx, client, "")

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SuccessWithQuery", func(t *testing.T) {
		mockRepo.On("List", ctx, client).Return(posts, nil).Once()

		result, err := service.List(ctx, client, "secret")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "2", result[0].ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.On("List", ctx, client).Return(nil, errors.New("list error")).Once()

		result, err := service.List(ctx, client, "")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestPostService_CRUD(t *testing.T) {
	mockRepo := new(repositories.MockPostRepository)
	service := NewPostService(mockRepo)
	ctx := context.Background()
	client := &pb.Client{}

	t.Run("GetSuccess", func(t *testing.T) {
		expected := &pb.Post{ID: "1", Title: "Title"}
		mockRepo.On("Get", ctx, client, "1").Return(expected, nil).Once()

		result, err := service.Get(ctx, client, "1")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("CreateSuccess", func(t *testing.T) {
		mockRepo.On("Create", ctx, client, "New", "Content", true).Return(nil).Once()

		err := service.Create(ctx, client, "New", "Content", true)

		assert.NoError(t, err)
	})

	t.Run("UpdateSuccess", func(t *testing.T) {
		id := "1"
		title := "Updated"
		content := "Updated Content"
		isPublic := false

		mockRepo.On("Update", ctx, client, id, map[string]any{
			"title":   title,
			"content": content,
			"public":  isPublic,
		}).Return(nil).Once()

		err := service.Update(ctx, client, id, title, content, isPublic)

		assert.NoError(t, err)
	})

	t.Run("DeleteSuccess", func(t *testing.T) {
		mockRepo.On("Delete", ctx, client, "1").Return(nil).Once()

		err := service.Delete(ctx, client, "1")

		assert.NoError(t, err)
	})
}
