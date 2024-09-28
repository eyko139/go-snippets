package main

import (
	"net/http"
	"time"

	"github.com/eyko139/go-snippets/config"
    "fmt"
)


func main() {

    env := config.NewEnv()

	cfg, err := config.NewApp(env)


	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	cfg.InfoLog.Printf("Starting server on %s", env.ServerPort) 
	srv := &http.Server{
        Addr:     fmt.Sprintf(":%s", env.ServerPort),
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
