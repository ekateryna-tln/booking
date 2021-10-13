package hendlers

import (
	"github.com/ekateryna-tln/booking/pkg/config"
	"github.com/ekateryna-tln/booking/pkg/models"
	"github.com/ekateryna-tln/booking/pkg/render"
	"net/http"
)

// Repo the repository used by handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	AppConfig *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		AppConfig: a,
	}
}

// SetHandlersRepo set the repository for the handlers
func SetHandlersRepo(r *Repository) {
	Repo = r
}

// Home is the about page handler
func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	rp.AppConfig.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Test test test"
	stringMap["remote_ip"] = rp.AppConfig.Session.GetString(r.Context(), "remote_ip")
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals renders the room page
func (rp *Repository) Generals(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (rp *Repository) Majors(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (rp *Repository) Availability(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "search-availability.page.tmpl", &models.TemplateData{})

}

// Reservation renders the make a reservation page and display form
func (rp *Repository) Reservation(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Contacts renders the contacts page and display form
func (rp *Repository) Contacts(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, "contacts.page.tmpl", &models.TemplateData{})
}
