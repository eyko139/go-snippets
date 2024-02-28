package util

import (
	"github.com/eyko139/go-snippets/internal/models"
)

type TemplateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
	FormErrors  map[string]string
}
