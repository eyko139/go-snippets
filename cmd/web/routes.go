package main

import (
	"github.com/eyko139/go-snippets/config"
	"net/http"
	"github.com/justinas/alice"
)

func Routes(cfg *config.Config) http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home(cfg))
	mux.HandleFunc("/snippet/view", snippetView(cfg))
	mux.HandleFunc("/snippet/create", snippetCreate(cfg))

	standard := alice.New(cfg.PanicRecovery, cfg.LogRequests, secureHeaders)
	return standard.Then(mux)
}
