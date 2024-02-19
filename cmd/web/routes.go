package main

import (
	"github.com/eyko139/go-snippets/config"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func Routes(cfg *config.Config) http.Handler {

	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Hlp.NotFound(w)
	})

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", home(cfg))
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", snippetView(cfg))
	router.HandlerFunc(http.MethodGet, "/snippet/create", snippetCreate(cfg))
	router.HandlerFunc(http.MethodPost, "/snippet/create", snippetCreatePost(cfg))

	standard := alice.New(cfg.PanicRecovery, cfg.LogRequests, secureHeaders)
	return standard.Then(router)
}
