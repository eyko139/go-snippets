package main

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"

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
		if err != nil {
			cfg.Hlp.ServerError(w, err)
			return
		}
		data := cfg.Hlp.NewTemplateData(r)
		data.Snippets = latest

		cfg.Hlp.Render(w, http.StatusOK, "home.html", data)
		if err != nil {
			cfg.Hlp.ServerError(w, err)
			return
		}
	}
}

func snippetCreate(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := cfg.Hlp.NewTemplateData(r)
		cfg.Hlp.Render(w, http.StatusOK, "create.html", data)
	}
}

func snippetView(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil || id < 1 {
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
		data := cfg.Hlp.NewTemplateData(r)
		data.Snippet = snippet

		cfg.Hlp.Render(w, http.StatusOK, "view.html", data)
		if err != nil {
			cfg.Hlp.ServerError(w, err)
			return
		}
	}
}

func snippetCreatePost(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			// headers are automatically converted into correct case according to http standard
			w.Header().Set("Allow", http.MethodPost)
			cfg.Hlp.ClientError(w, http.StatusMethodNotAllowed)
			return
		}
		err := r.ParseForm()
		if err != nil {
			cfg.Hlp.ClientError(w, http.StatusBadRequest)
			return
		}
		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")
		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			cfg.Hlp.ClientError(w, http.StatusBadRequest)
			return
		}
		//  Validate user input
		fieldErrors := make(map[string]string)
		if strings.TrimSpace(title) == "" {
			fieldErrors["title"] = "This field cannot be blank"
		} else if utf8.RuneCountInString(title) > 100 {
			fieldErrors["title"] = "Field cannot be longer than 100 characters"
		}

		if strings.TrimSpace(content) == "" {
			fieldErrors["content"] = "Content cannot be empty"
		}

		expirationDates := []int{1, 7, 365}
		if !slices.Contains(expirationDates, expires) {
			fieldErrors["expires"] = fmt.Sprintf("Expiration date may only contain one of the following values %d", expires)
		}

		id, err := cfg.Snippets.Insert(title, content, expires)
		if err != nil {
			cfg.Hlp.ServerError(w, err)
		}
		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id),
			http.StatusSeeOther)
	}
}
