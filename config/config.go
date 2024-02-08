package config

import (
	"database/sql"
	"github.com/eyko139/go-snippets/cmd/util"
	"github.com/eyko139/go-snippets/internal/models"
	"html/template"
	"log"
	"os"
)

type Config struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Hlp           *util.Helpers
	Snippets      *models.SnippetModel
	TemplateCache map[string]*template.Template
}

func New(db *sql.DB) (*Config, error) {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	tc, err := models.NewTemplateCache()
	if err != nil {
		return nil, err
	}
	return &Config{
		ErrorLog: errLog,
		InfoLog:  infoLog,
		Hlp:      util.NewHelper(tc, errLog, infoLog),
		Snippets: &models.SnippetModel{DB: db},
	}, nil
}
