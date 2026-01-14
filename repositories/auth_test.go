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

func TestPBAuthRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("Authenticate_Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/collections/users/auth-with-password", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"token": "test-token",
				"record": map[string]any{
					"id":    "user-1",
					"email": "test@example.com",
					"name":  "Test User",
				},
			})
		}))
		defer server.Close()

		client := pb.NewClient(server.URL)
		repo := NewAuthRepository()

		token, user, err := repo.Authenticate(ctx, client, "test@example.com", "password")

		assert.NoError(t, err)
		assert.Equal(t, "test-token", token)
		assert.Equal(t, "user-1", user.ID)
	})

	t.Run("Authenticate_Failure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		client := pb.NewClient(server.URL)
		repo := NewAuthRepository()

		token, user, err := repo.Authenticate(ctx, client, "test@example.com", "wrong")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
	})

	t.Run("CreateUser_Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/collections/users/records", r.URL.Path)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{"id": "user-2"})
		}))
		defer server.Close()

		client := pb.NewClient(server.URL)
		repo := NewAuthRepository()

		err := repo.CreateUser(ctx, client, map[string]any{"email": "new@example.com"})

		assert.NoError(t, err)
	})
}
