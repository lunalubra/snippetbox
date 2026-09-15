package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/lunalubra/snippetbox/assets"
	"github.com/lunalubra/snippetbox/internal/models"
)

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// canExtend reports whether a snippet with the given expiry time is close
// enough to expiring for its expiry to be extended.
func canExtend(expires time.Time) bool {
	now := time.Now()

	return expires.After(now) && expires.Sub(now) <= models.ExtensionWindow
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"canExtend": canExtend,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(assets.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(assets.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}
