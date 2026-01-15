package views

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexView(t *testing.T) {
	t.Run("LoggedIn", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := Index(true).Render(context.Background(), buf)
		assert.NoError(t, err)

		content := buf.String()
		assert.Contains(t, content, "Launch Your")
		assert.Contains(t, content, "Full-Stack Mission")
		assert.NotContains(t, content, "Ready to Launch?") // CTA should be hidden
	})

	t.Run("LoggedOut", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := Index(false).Render(context.Background(), buf)
		assert.NoError(t, err)

		content := buf.String()
		assert.Contains(t, content, "Launch Your")
		assert.Contains(t, content, "Full-Stack Mission")
		assert.Contains(t, content, "Ready to Launch?") // CTA should be visible
	})
}

func TestDashboardView(t *testing.T) {
	userName := "Commander Shepard"
	userEmail := "shepard@normandy.sr2"
	postCount := 42
	csrf := "fake-csrf-token"

	buf := new(bytes.Buffer)
	err := Dashboard(userName, userEmail, postCount, csrf).Render(context.Background(), buf)
	assert.NoError(t, err)

	content := buf.String()
	assert.Contains(t, content, "Command Center")
	assert.Contains(t, content, "Welcome aboard")
	assert.Contains(t, content, userName)
	assert.Contains(t, content, userEmail)
	assert.Contains(t, content, "42")
	assert.Contains(t, content, csrf)
	assert.Contains(t, content, "name=\"_csrf\"")
	assert.Contains(t, content, "EXECUTE_TRANSMISSION") // Verify the new HTMX/styled button
}
