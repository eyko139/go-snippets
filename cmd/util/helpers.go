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

func NewHelper(templateCache map[string]*template.Template, err *log.Logger, info *log.Logger) *Helpers {
	return &Helpers{
		ErrorLog:      err,
		InfoLog:       info,
		TemplateCache: templateCache,
	}
}

func (h *Helpers) NewTemplateData() *TemplateData{
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
