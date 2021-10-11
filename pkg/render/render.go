package render

import (
	"github.com/ekateryna-tln/booking/pkg/config"
	"github.com/ekateryna-tln/booking/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = make(template.FuncMap)

var appConfig *config.AppConfig

// SetRenderAppConfig sets the config for the render template package
func SetRenderAppConfig(a *config.AppConfig) {
	appConfig = a
}

func AddDefaultData(tmplData *models.TemplateData) *models.TemplateData {
	return tmplData
}

// RenderTemplate
func RenderTemplate(w http.ResponseWriter, tmpl string, tmplData *models.TemplateData) {
	var templateList map[string]*template.Template
	if appConfig.UseCache {
		templateList = appConfig.TemplateCache
	} else {
		templateList = GetTemplateCache()
	}

	template, exist := templateList[tmpl]
	if !exist {
		log.Fatal("could not get template from template cache")
	}

	tmplData = AddDefaultData(tmplData)
	err := template.Execute(w, tmplData)
	if err != nil {
		log.Fatal("writing template error:", err)
	}
}

func GetTemplateCache() map[string]*template.Template {
	templateCache, err := CreateTemplateCache()
	if err != nil {
		log.Fatal("create template cache error:", err)
	}
	return templateCache
}

// CreateTemplateCache create a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return templateCache, err
	}

	layouts, err := filepath.Glob("./templates/*.layout.tmpl")
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
			template, err = template.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = template
	}
	return templateCache, nil
}
