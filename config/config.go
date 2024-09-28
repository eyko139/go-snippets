package config

import (
	"context"
	"html/template"
	"log"
	"os"

	"github.com/eyko139/go-snippets/internal/session"
	"github.com/eyko139/go-snippets/internal/session/providers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	UserModel      models.UserModelInterface
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func NewApp(appConfig *Env) (*Config, error) {

	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(appConfig.DBConnectionString))

    if err != nil {
        errLog.Printf("Failed to connect to DB %s", err)
    }
    pingErr := client.Ping(context.Background(), nil)

    if pingErr != nil {
        errLog.Printf("Database Ping failed", pingErr)
    }

	providers.InitSessionProvider(client)
	globalSessions, err := session.NewManager(appConfig.SessionProvider, "gosessionid", 360)

	if err != nil {
		panic("Could not initialize session manager")
	}

	go globalSessions.GC()

	tc, err := models.NewTemplateCache()

	if err != nil {
		return nil, err
	}

	return &Config{
		ErrorLog:       errLog,
		InfoLog:        infoLog,
		Hlp:            util.NewHelper(tc, errLog, infoLog),
		Snippets:       &models.SnippetModel{DBMongo: client},
		GlobalSessions: globalSessions,
		UserModel:      &models.UserModel{DbClient: client},
	}, nil
}
