package util

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type Helpers struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	TemplateCache map[string]*template.Template
}

type ErrorData struct {
	Title   string
	Message string
}

func NewHelper(templateCache map[string]*template.Template, err *log.Logger, info *log.Logger) *Helpers {
	return &Helpers{
		ErrorLog:      err,
		InfoLog:       info,
		TemplateCache: templateCache,
	}
}

func (h *Helpers) NewTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		CurrentYear: time.Now().Year(),
	}
}

func (h *Helpers) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// set frame depth to 2, we don't want to see this line first on the stack trace
	// when error occurs
	h.ErrorLog.Output(2, trace)
	h.ErrorLog.Print(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *Helpers) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (h *Helpers) NotFound(w http.ResponseWriter) {
	h.ClientError(w, http.StatusNotFound)
}
func (h *Helpers) ReturnTemplateError(w http.ResponseWriter, data map[string]string) error {
	// Create a new TemplateData struct holding the error map
	templateData := &TemplateData{
		FormErrors: data,
	}

	buf := new(bytes.Buffer)
	ts := h.TemplateCache["create.html"]
	err := ts.ExecuteTemplate(buf, "base", templateData)
	if err != nil {
		h.ServerError(w, err)
	}
	w.Write([]byte(buf.String()))
	return nil
}

func (h *Helpers) Render(w http.ResponseWriter, status int, page string, data *TemplateData) {
	ts, ok := h.TemplateCache[page]

	// Write the template to a buffer first to catch runtime errors
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	if !ok {
		err := fmt.Errorf("Template %s doesn't exist", page)
		h.ServerError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)

	if err != nil {
		h.ServerError(w, err)
	}
}
