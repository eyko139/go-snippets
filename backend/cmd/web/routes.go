package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (cfg *Config) Routes() http.Handler {

	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./backend/ui/static/"))

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Hlp.NotFound(w)
	})

	dynamic := alice.New(cfg.SessionContext)

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.Handler(http.MethodGet, "/health", dynamic.ThenFunc(health()))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(cfg.home()))

	router.Handler(http.MethodGet, "/ws", cfg.ws())
	router.Handler(http.MethodGet, "/wsinfo", cfg.wsInfo())

	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(cfg.snippetView()))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(cfg.snippetCreate()))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(cfg.snippetCreatePost()))
	router.Handler(http.MethodPost, "/temp", dynamic.ThenFunc(cfg.tempContentPost()))
	router.Handler(http.MethodGet, "/snippets", dynamic.ThenFunc(cfg.getSnippets()))

	// User Management
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(cfg.userSignup()))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(cfg.userSignupPost()))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(cfg.userLogin()))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(cfg.userLoginPost()))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(cfg.userLogoutPost()))

	standard := alice.New(cfg.PanicRecovery, cfg.LogRequests, secureHeaders)
	return standard.Then(router)
}
