package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/eyko139/go-snippets/cmd/util"
	"github.com/eyko139/go-snippets/internal/models"
	"github.com/eyko139/go-snippets/internal/session"
	"github.com/eyko139/go-snippets/internal/validator"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (cfg *Config) home() http.HandlerFunc {
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
		data := cfg.Hlp.NewTemplateData(r)
		data.Snippets = latest

		cfg.Hlp.Render(w, http.StatusOK, "home.html", data)
	}
}

func (cfg *Config) snippetCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := cfg.Hlp.NewTemplateData(r)
		ctx := r.Context()
		content := ctx.Value("session").(session.Session).Get("content")
		if content != nil {
			data.Content = content.(string)
		} else {
			data.Content = ""
		}
		cfg.Hlp.Render(w, http.StatusOK, "create.html", data)
	}
}

func (cfg *Config) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := cfg.Hlp.NewTemplateData(r)
		cfg.Hlp.Render(w, http.StatusOK, "login.html", data)
	}
}

func (cfg *Config) loginPost() http.HandlerFunc {
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

func (cfg *Config) snippetView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		id := params.ByName("id")
		if id == "" {
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
		data := cfg.Hlp.NewTemplateData(r)
		data.Snippet = snippet

		cfg.Hlp.Render(w, http.StatusOK, "view.html", data)
	}
}

// Save the temporay input of the the snippet textarea
// and store in the session, this tempContent will be retrieved when visiting
// the createSnippet page
func (cfg *Config) tempContentPost() http.HandlerFunc {
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

func (cfg *Config) getTempContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content := cfg.GlobalSessions.SessionStart(w, r).Get("content")
		w.Write([]byte(content.(string)))
	}
}

func (cfg *Config) snippetCreatePost() http.HandlerFunc {
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
		form := util.SnippetCreateForm{
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
		cfg.GlobalSessions.SessionStart(w, r).Set("flash", "snippped successfully created")
		cfg.Broker.Publish(content)
		cfg.GlobalSessions.SessionStart(w, r).Delete("content")
		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%s", id),
			http.StatusSeeOther)
	}
}

// User Management

func (cfg *Config) userSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := cfg.Hlp.NewTemplateData(r)
		data.Form = &util.UserSignupForm{}
		cfg.Hlp.Render(w, http.StatusOK, "signup.html", data)
	}
}

func (cfg *Config) userSignupPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			cfg.Hlp.ServerError(w, err)
		}
		formData := &util.UserSignupForm{
			Email:    r.PostForm.Get("email"),
			Name:     r.PostForm.Get("name"),
			Password: r.PostForm.Get("password"),
		}
		formData.CheckField(validator.NotBlank(formData.Name), "name", "This field cannot be blank")
		formData.CheckField(validator.NotBlank(formData.Email), "email", "This field cannot be blank")
		formData.CheckField(validator.Matches(formData.Email, validator.EmailRX), "email", "Please provide a valid email address")
		formData.CheckField(validator.NotBlank(formData.Password), "password", "Cannot be Blank")
		if !formData.Valid() {
			td := cfg.Hlp.NewTemplateData(r)
			td.FormErrors = formData.FieldErrors
			td.Form = formData
			cfg.Hlp.Render(w, http.StatusUnprocessableEntity, "signup.html", td)
			return
		}
		err = cfg.UserModel.Insert(formData.Name, formData.Email, formData.Password)

		if err != nil {
			cfg.Hlp.ServerError(w, err)
		}

		return
	}
}

func (cfg *Config) userLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	}
}
func (cfg *Config) userLoginPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	}
}

func (cfg *Config) userLogoutPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	}
}

func (cfg *Config) getSnippets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		latest, err := cfg.Snippets.Latest()
		if err != nil {
			http.Error(w, "error", 304)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(latest)
	}
}

func health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
