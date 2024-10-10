package mocks

import (
	"github.com/eyko139/go-snippets/internal/models"
	"time"
)

var mockSnippet = &models.Snippet{
	ID:      "1",
	Title:   "mockTitle",
	Content: "mockContent",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (string, error) {
	return "2", nil
}

func (m *SnippetModel) Get(id string) (*models.Snippet, error) {
	if id == "1" {
		return mockSnippet, nil
	}
	return nil, models.ErrNoRecord
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
