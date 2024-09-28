package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/eyko139/go-snippets/internal/session"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/eyko139/go-snippets/config"
	"github.com/eyko139/go-snippets/internal/session/providers"
)


func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:password@localhost:27017"))

	providers.InitSessionProvider(client)
	globalSessions, err := session.NewManager("mongo", "gosessionid", 360)

	if err != nil {
		panic("Could not initialize session manager")
	}

	go globalSessions.GC()



	addr := flag.String("addr", ":4000", "Http network address")
	flag.Parse()

	cfg, err := config.New(client, globalSessions)

	if err != nil {
		panic("Error creating config")
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	cfg.InfoLog.Printf("Starting server on %s", *addr)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: cfg.ErrorLog,
		Handler:  Routes(cfg),
		//NOTE: Always set IdleTimeout explicitly, otherwise IdleTimout = ReadTimeout
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = srv.ListenAndServe()
	cfg.ErrorLog.Fatal(err)
}
