package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eyko139/go-snippets/config"
	"github.com/eyko139/go-snippets/internal/models"
)

func home(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// http.NotFound helper for writing 404 and message
			cfg.Hlp.NotFound(w)
			return
		}
		latest, err := cfg.Snippets.Latest()
		if err != nil {
			cfg.Hlp.ServerError(w, err)
			return
		}
		//
		if err != nil {
			cfg.Hlp.ServerError(w, err)
			return
		}
		data := cfg.Hlp.NewTemplateData()
		data.Snippets = latest

		cfg.Hlp.Render(w, http.StatusOK, "home.html", data)
		if err != nil {
			cfg.Hlp.ServerError(w, err)
			return
		}
	}
}

func snippetView(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			cfg.Hlp.NotFound(w)
		}
		snippet, err := cfg.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				cfg.Hlp.NotFound(w)
			} else {
				cfg.Hlp.ServerError(w, err)
			}
			return
		}
		if err != nil {
			cfg.Hlp.ServerError(w, err)
			return
		}
		data := cfg.Hlp.NewTemplateData()
		data.Snippet = snippet

		cfg.Hlp.Render(w, http.StatusOK, "view.html", data)
		if err != nil {
			cfg.Hlp.ServerError(w, err)
			return
		}
	}
}

func snippetCreate(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			// headers are automatically converted into correct case according to http standard
			w.Header().Set("Allow", http.MethodPost)
			cfg.Hlp.ClientError(w, http.StatusMethodNotAllowed)
			return
		}
		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n Kobayashi Issa"
		expires := 7
		id, err := cfg.Snippets.Insert(title, content, expires)
		if err != nil {
			cfg.Hlp.ServerError(w, err)
		}
		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id),
			http.StatusSeeOther)
		w.Write([]byte("Creating new snippet"))
	}
}
