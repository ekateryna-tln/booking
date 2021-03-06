package render

import (
	"errors"
	"fmt"
	"github.com/ekateryna-tln/booking/internal/config"
	"github.com/ekateryna-tln/booking/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

var functions = template.FuncMap{
	"formatDate":       FormatDate,
	"formatDateString": FormatDateString,
	"iterate":          Iterate,
	"add":              Add,
}

var app *config.App
var pathToTemplates = "./templates"

// SetRenderApp sets the config for the render template package
func SetRenderApp(a *config.App) {
	app = a
}

// FormatDate returns date in YYYY-MM-DD format
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateString returns date in YYYY-MM-DD format
func FormatDateString(t time.Time, f string) string {
	return t.Format(f)
}

func Add(a, b int) int {
	return a + b
}

// Iterate returns a slice of ints, starting at 1, going to count
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

func AddDefaultData(tmplData *models.TemplateData, r *http.Request) *models.TemplateData {
	tmplData.Flash = app.Session.PopString(r.Context(), "flash")
	tmplData.Error = app.Session.PopString(r.Context(), "error")
	tmplData.Warning = app.Session.PopString(r.Context(), "warning")
	tmplData.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		tmplData.IsAuthenticated = 1
	}
	return tmplData
}

// Template renders templates using http/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, tmplData *models.TemplateData) error {
	var templateList map[string]*template.Template
	if app.UseCache {
		templateList = app.TemplateCache
	} else {
		var err error
		templateList, err = CreateTemplateCache()
		if err != nil {
			return err
		}
	}

	template, ok := templateList[tmpl]
	if !ok {
		return errors.New("could not get template from template cache")
	}

	tmplData = AddDefaultData(tmplData, r)
	err := template.Execute(w, tmplData)
	if err != nil {
		return err
	}
	return nil
}

// CreateTemplateCache create a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := make(map[string]*template.Template)

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return templateCache, err
	}

	layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		template, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		if len(layouts) > 0 {
			template, err = template.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = template
	}
	return templateCache, nil
}
