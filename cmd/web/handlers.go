package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/eyko139/go-snippets/config"
	"github.com/eyko139/go-snippets/internal/models"
	"github.com/eyko139/go-snippets/internal/validator"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

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
		content := cfg.GlobalSessions.SessionStart(w, r).Get("content")
		if content != nil {
			//data.Content = content.(string)
		} else {
			data.Content = ""
		}
		cfg.Hlp.Render(w, http.StatusOK, "create.html", data)
	}
}

func login(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := cfg.Hlp.NewTemplateData(r)
		cfg.Hlp.Render(w, http.StatusOK, "login.html", data)
	}
}

func loginPost(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := cfg.GlobalSessions.SessionStart(w, r)

		err := r.ParseForm()
		if err != nil {
			cfg.Hlp.ClientError(w, http.StatusBadRequest)
			return
		}
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")
		sess.Set("username", username)
		data := cfg.Hlp.NewTemplateData(r)
		fmt.Println(username, password)
		cfg.Hlp.Render(w, http.StatusOK, "login.html", data)
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

func tempContentPost(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			cfg.Hlp.ClientError(w, http.StatusBadRequest)
			return
		}
		content := r.PostForm.Get("content")
		cfg.GlobalSessions.SessionStart(w, r).Set("content", content)
	}
}

func getTempContent(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content := cfg.GlobalSessions.SessionStart(w, r).Get("content")
		w.Write([]byte(content.(string)))
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
		form := snippetCreateForm{
			Title:   title,
			Content: content,
			Expires: expires,
		}

		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
		form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "content", "This field cannot be blank")

		if !form.Valid() {
			data := cfg.Hlp.NewTemplateData(r)
			data.FormErrors = form.FieldErrors
			content := cfg.GlobalSessions.SessionStart(w, r).Get("content")
			if content != nil {
				data.Content = content.(string)
			} else {
				data.Content = ""
			}
			err := cfg.Hlp.ReturnTemplateError(w, data)
			if err != nil {
				cfg.Hlp.ServerError(w, err)
			}
			return
		}

		id, err := cfg.Snippets.Insert(title, content, expires)
		if err != nil {
			cfg.Hlp.ServerError(w, err)
		}
		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id),
			http.StatusSeeOther)
	}
}
