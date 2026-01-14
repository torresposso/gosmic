package pb

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthWithPassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/collections/users/auth-with-password", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "test@example.com", body["identity"])
		assert.Equal(t, "password123", body["password"])

		resp := authResponse{
			Token: "test-token",
			Record: User{
				ID:    "user-123",
				Email: "test@example.com",
				Name:  "Test User",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	token, user, err := client.AuthWithPassword("test@example.com", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "test-token", token)
	assert.Equal(t, "user-123", user.ID)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestListPosts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/collections/posts/records", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		resp := listPostsResponse{
			Items: []Post{
				{ID: "1", Title: "Post 1"},
				{ID: "2", Title: "Post 2"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	posts, err := client.ListPosts()

	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, "Post 1", posts[0].Title)
}

func TestCreatePost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/collections/posts/records", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "New Post", body["title"])
		assert.Equal(t, "Content", body["content"])

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithToken("test-token")
	err := client.CreatePost("New Post", "Content", true)

	assert.NoError(t, err)
}
