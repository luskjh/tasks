package lib

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

type Registry map[string]*template.Template

func NewRegistry(fsys fs.FS) (Registry, error) {
	var pages, rest []string

	if err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		if strings.HasPrefix(path, "page/") {
			pages = append(pages, path)
		} else {
			rest = append(rest, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	base := template.New("")
	for _, path := range rest {
		if err := parse(fsys, base, path); err != nil {
			return nil, err
		}
	}

	registry := make(Registry, len(pages))
	for _, page := range pages {
		b, _ := base.Clone()
		if err := parse(fsys, b, page); err != nil {
			return nil, err
		}

		key := strings.TrimPrefix(page, "page/")
		registry[key] = b
	}

	return registry, nil
}

func (r Registry) RenderWithLayout(w http.ResponseWriter, page string, data any, layout string) {
	t, ok := r[page]
	if !ok {
		// TODO: logging
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(w, layout, data); err != nil {
		// TODO: logging
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func parse(fsys fs.FS, t *template.Template, path string) error {
	b, err := fs.ReadFile(fsys, path)
	if err != nil {
		return err
	}

	tmpl := t
	if tmpl.Name() != path {
		tmpl = t.New(path)
	}

	_, err = tmpl.Parse(string(b))
	return err
}
