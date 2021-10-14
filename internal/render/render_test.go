package render

import (
	"github.com/ekateryna-tln/booking/internal/models"
	"net/http"
	"testing"
)

func TestSetRenderApp(t *testing.T) {
	SetRenderApp(app)
}

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "success flash")
	result := AddDefaultData(&td, r)
	if result.Flash != "success flash" {
		t.Error("the value \"success flash\" not found in session")
	}
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	var ww myWriter

	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = RenderTemplate(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("get", "/", nil)
	if err != nil {
		return nil, err
	}
	context := r.Context()
	context, _ = session.Load(context, r.Header.Get("X-Session"))
	r = r.WithContext(context)
	return r, nil
}
