package config

import (
	"database/sql"
	"github.com/eyko139/go-snippets/internal/session"
	"html/template"
	"log"
	"os"

	"github.com/eyko139/go-snippets/cmd/util"
	"github.com/eyko139/go-snippets/internal/models"
)

type Config struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Hlp            *util.Helpers
	Snippets       *models.SnippetModel
	TemplateCache  map[string]*template.Template
	GlobalSessions *session.Manager
}

func New(db *sql.DB, manager *session.Manager) (*Config, error) {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	tc, err := models.NewTemplateCache()
	if err != nil {
		return nil, err
	}
	return &Config{
		ErrorLog:       errLog,
		InfoLog:        infoLog,
		Hlp:            util.NewHelper(tc, errLog, infoLog),
		Snippets:       &models.SnippetModel{DB: db},
		GlobalSessions: manager,
	}, nil
}
