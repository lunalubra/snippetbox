package mocks

import (
	"time"

	"github.com/lunalubra/snippetbox/internal/models"
)

var mockSnippet = models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now().Add(24 * time.Hour),
}

var mockLongLivedSnippet = models.Snippet{
	ID:      3,
	Title:   "Over the wintry forest",
	Content: "Over the wintry forest...",
	Created: time.Now(),
	Expires: time.Now().Add(100 * 24 * time.Hour),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	case 3:
		return mockLongLivedSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}

func (m *SnippetModel) ExtendExpiry(id int) (time.Time, error) {
	switch id {
	case 1:
		return mockSnippet.Expires.Add(models.ExtensionPeriod), nil
	default:
		return time.Time{}, models.ErrNotExtendable
	}
}
