package hendlers

import (
	"encoding/json"
	"fmt"
	"github.com/ekateryna-tln/booking/internal/config"
	"github.com/ekateryna-tln/booking/internal/models"
	"github.com/ekateryna-tln/booking/internal/render"
	"log"
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
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Test test test"
	stringMap["remote_ip"] = rp.AppConfig.Session.GetString(r.Context(), "remote_ip")
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals renders the room page
func (rp *Repository) Generals(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, request, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (rp *Repository) Majors(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, request, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (rp *Repository) Availability(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, request, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (rp *Repository) PostAvailability(writer http.ResponseWriter, request *http.Request) {
	startDate := request.Form.Get("start_date")
	endDate := request.Form.Get("end_date")
	writer.Write([]byte(fmt.Sprintf("Start date is %s, end date is %s", startDate, endDate)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (rp *Repository) AvailabilityJSON(writer http.ResponseWriter, request *http.Request) {
	response := jsonResponse{
		OK:      true,
		Message: "Available!",
	}
	out, err := json.MarshalIndent(response, "", "     ")
	if err != nil {
		log.Fatalln(err)
	}
	writer.Header().Set("Content-type", "application/json")
	writer.Write(out)
}

// Reservation renders the make a reservation page and display form
func (rp *Repository) Reservation(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, request, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Contacts renders the contacts page and display form
func (rp *Repository) Contacts(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, request, "contacts.page.tmpl", &models.TemplateData{})
}
