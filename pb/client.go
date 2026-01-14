package pb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	AuthToken  string
	HTTPClient *http.Client
	AuthRecord *User
}

func (c *Client) newRequest(method, path string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	return req, nil
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type authResponse struct {
	Token  string `json:"token"`
	Record User   `json:"record"`
}

type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Public  bool   `json:"public"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type listPostsResponse struct {
	Items []Post `json:"items"`
}

func NewClient(url string) *Client {
	return &Client{
		BaseURL: url,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// WithToken creates a request-scoped client with a specific user's token.
// It shares the underlying HTTP connection pool for efficiency.
func (c *Client) WithToken(token string) *Client {
	return &Client{
		BaseURL:    c.BaseURL,
		HTTPClient: c.HTTPClient, // Shared connection pool
		AuthToken:  token,
	}
}

func (c *Client) IsAuthenticated() bool {
	return c.AuthToken != ""
}

func (c *Client) GetCurrentUserEmail() string {
	if c.AuthRecord != nil {
		return c.AuthRecord.Email
	}
	return ""
}

func (c *Client) GetCurrentUserName() string {
	if c.AuthRecord != nil {
		if c.AuthRecord.Name != "" {
			return c.AuthRecord.Name
		}
		return c.AuthRecord.Email
	}
	return "User"
}

func (c *Client) GetUserID() string {
	if c.AuthRecord != nil {
		return c.AuthRecord.ID
	}
	return ""
}

// AuthWithPassword authenticates and returns (token, *User, error).
// The token should be stored in a cookie by the caller.
func (c *Client) AuthWithPassword(email, password string) (string, *User, error) {
	body := map[string]any{
		"identity": email,
		"password": password,
	}

	req, err := c.newRequest("POST", "/api/collections/users/auth-with-password", body)
	if err != nil {
		return "", nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, errors.New("authentication failed")
	}

	var authResp authResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return authResp.Token, &authResp.Record, nil
}

func (c *Client) Logout() {
	c.AuthToken = ""
	c.AuthRecord = nil
}

func (c *Client) ListPosts() ([]Post, error) {
	req, err := c.newRequest("GET", "/api/collections/posts/records", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch posts: %d", resp.StatusCode)
	}

	var listResp listPostsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Items, nil
}

func (c *Client) GetPost(id string) (*Post, error) {
	req, err := c.newRequest("GET", "/api/collections/posts/records/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch post: %d", resp.StatusCode)
	}

	var post Post
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &post, nil
}

func (c *Client) CreatePost(title, content string, isPublic bool) error {
	if c.AuthToken == "" {
		return errors.New("unauthorized")
	}

	body := map[string]interface{}{
		"title":   title,
		"content": content,
		"author":  c.GetUserID(),
		"public":  isPublic,
	}

	req, err := c.newRequest("POST", "/api/collections/posts/records", body)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create post: %s", string(respBody))
	}

	return nil
}

func (c *Client) DeletePost(id string) error {
	if c.AuthToken == "" {
		return errors.New("unauthorized")
	}

	req, err := c.newRequest("DELETE", "/api/collections/posts/records/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete post: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) UpdatePost(id string, data map[string]any) error {
	if c.AuthToken == "" {
		return errors.New("unauthorized")
	}

	req, err := c.newRequest("PATCH", "/api/collections/posts/records/"+id, data)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update post: %d", resp.StatusCode)
	}
	return nil
}

// CreateRecord generic helper for registration
func (c *Client) CreateRecord(collection string, data map[string]any) error {
	req, err := c.newRequest("POST", "/api/collections/"+collection+"/records", data)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create record: %s", string(respBody))
	}

	return nil
}
