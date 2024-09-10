package config

import (
	"html/template"
	"log"
	"os"

	"github.com/eyko139/go-snippets/internal/session"
	"go.mongodb.org/mongo-driver/mongo"

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
    UserModel models.UserModelInterface
}

func New(mongoClient *mongo.Client, manager *session.Manager) (*Config, error) {
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
		Snippets:       &models.SnippetModel{DBMongo: mongoClient},
		GlobalSessions: manager,
        UserModel: &models.UserModel{DbClient: mongoClient},
	}, nil
}
