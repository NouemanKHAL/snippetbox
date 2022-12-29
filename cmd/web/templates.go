package main

import (
	"html/template"
	"path/filepath"

	"github.com/NouemanKHAL/snippetbox/internal/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./ui/html/pages/*.go.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// parse the base template
		ts, err := template.ParseFiles("./ui/html/base.go.tmpl")
		if err != nil {
			return nil, err
		}

		// parse the partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.go.tmpl")
		if err != nil {
			return nil, err
		}

		// parse the page
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
