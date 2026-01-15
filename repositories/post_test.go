package repositories

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/torresposso/gosmic/pb"
)

func TestPBPostRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("List_Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/collections/posts/records", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"items": []map[string]any{
					{"id": "p1", "title": "Post 1"},
				},
			})
		}))
		defer server.Close()

		client := pb.NewClient(server.URL)
		repo := NewPostRepository()

		posts, err := repo.List(ctx, client)
		assert.NoError(t, err)
		assert.Len(t, posts, 1)
		assert.Equal(t, "p1", posts[0].ID)
	})

	t.Run("Get_Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/collections/posts/records/p1", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{"id": "p1", "title": "Post 1"})
		}))
		defer server.Close()

		client := pb.NewClient(server.URL)
		repo := NewPostRepository()

		post, err := repo.Get(ctx, client, "p1")
		assert.NoError(t, err)
		assert.Equal(t, "Post 1", post.Title)
	})

	t.Run("Create_Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/collections/posts/records", r.URL.Path)
			w.WriteHeader(http.StatusCreated)
		}))
		defer server.Close()

		client := pb.NewClient(server.URL)
		client.AuthToken = "dummy" // Required by client.CreatePost
		repo := NewPostRepository()

		err := repo.Create(ctx, client, "New", "Content", true)
		assert.NoError(t, err)
	})

	t.Run("Update_Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/collections/posts/records/p1", r.URL.Path)
			assert.Equal(t, "PATCH", r.Method)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := pb.NewClient(server.URL)
		client.AuthToken = "dummy"
		repo := NewPostRepository()

		err := repo.Update(ctx, client, "p1", map[string]any{"title": "Updated"})
		assert.NoError(t, err)
	})

	t.Run("Delete_Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/collections/posts/records/p1", r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client := pb.NewClient(server.URL)
		client.AuthToken = "dummy"
		repo := NewPostRepository()

		err := repo.Delete(ctx, client, "p1")
		assert.NoError(t, err)
	})

	t.Run("TogglePublic_Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				assert.Equal(t, "/api/collections/posts/records/p1", r.URL.Path)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{"id": "p1", "public": false})
			} else if r.Method == "PATCH" {
				assert.Equal(t, "/api/collections/posts/records/p1", r.URL.Path)
				var body map[string]any
				json.NewDecoder(r.Body).Decode(&body)
				assert.True(t, body["public"].(bool))
				w.WriteHeader(http.StatusOK)
			}
		}))
		defer server.Close()

		client := pb.NewClient(server.URL)
		client.AuthToken = "dummy"
		repo := NewPostRepository()

		err := repo.TogglePublic(ctx, client, "p1")
		assert.NoError(t, err)
	})
}
