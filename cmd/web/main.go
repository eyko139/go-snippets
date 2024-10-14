package main

import (
	"fmt"
	"github.com/eyko139/go-snippets/cmd/web/websocket"
	"net/http"
	"time"
)

func main() {

	env := NewEnv()

	hub := websocket.NewHub()
	go hub.Run()

	cfg, err := NewApp(env, hub)

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	cfg.InfoLog.Printf("Starting server on %s", env.ServerPort)
	srv := &http.Server{
		Addr:     fmt.Sprintf(":%s", env.ServerPort),
		ErrorLog: cfg.ErrorLog,
		Handler:  cfg.Routes(),
		//NOTE: Always set IdleTimeout explicitly, otherwise IdleTimout = ReadTimeout
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = srv.ListenAndServe()
	cfg.ErrorLog.Fatal(err)
}
