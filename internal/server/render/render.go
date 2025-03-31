package render

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type Renderer struct {
	templates map[string]*template.Template
	mu        sync.RWMutex
	fs        fs.FS
}

func New(fs fs.FS) *Renderer {
	return &Renderer{
		templates: make(map[string]*template.Template),
		fs:        fs,
	}
}

func (r *Renderer) LoadTemplates() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	basePath := "templates"
	layoutPath := filepath.Join(basePath, "layouts", "base.html")

	return fs.WalkDir(r.fs, basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-html files
		if d.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		// Skip layout files
		if filepath.Dir(path) == filepath.Join(basePath, "layouts") {
			return nil
		}

		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}

		tmpl := template.New(filepath.Base(path)).Funcs(template.FuncMap{
			"isActiveNav": func(current, target string) bool {
				return current == target
			},
		})

		// Parse layout first
		tmpl, err = tmpl.ParseFS(r.fs, layoutPath)
		if err != nil {
			return err
		}

		// Parse the specific template
		tmpl, err = tmpl.ParseFS(r.fs, path)
		if err != nil {
			return err
		}

		r.templates[relPath] = tmpl
		return nil
	})
}

func (r *Renderer) HTML(w http.ResponseWriter, status int, name string, data interface{}) {
	r.mu.RLock()
	tmpl, ok := r.templates[name]
	r.mu.RUnlock()

	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering template %s: %v", name, err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
