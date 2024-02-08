package models

import (
	"html/template"
	"path/filepath"
	"time"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:12")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("/home/luk/GolandProjects/snippet/ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("/home/luk/GolandProjects/snippet/ui/html/base.html")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("/home/luk/GolandProjects/snippet/ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)

		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
