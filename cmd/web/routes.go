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

    dynamic := alice.New(cfg.SessionContext)

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(home(cfg)))
    router.Handler(http.MethodGet, "/login", dynamic.ThenFunc(login(cfg)))
    router.Handler(http.MethodPost, "/login", dynamic.ThenFunc(loginPost(cfg)))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(snippetView(cfg)))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(snippetCreate(cfg)))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(snippetCreatePost(cfg)))
	router.Handler(http.MethodPost, "/temp", dynamic.ThenFunc(tempContentPost(cfg)))

	standard := alice.New(cfg.PanicRecovery, cfg.LogRequests, secureHeaders)
	return standard.Then(router)
}
