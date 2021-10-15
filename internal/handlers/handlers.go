package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ekateryna-tln/booking/internal/config"
	"github.com/ekateryna-tln/booking/internal/forms"
	"github.com/ekateryna-tln/booking/internal/helpers"
	"github.com/ekateryna-tln/booking/internal/models"
	"github.com/ekateryna-tln/booking/internal/render"
	"net/http"
)

// Repo the repository used by handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.App
}

// NewRepo creates a new repository
func NewRepo(a *config.App) *Repository {
	return &Repository{
		App: a,
	}
}

// SetHandlersRepo set the repository for the handlers
func SetHandlersRepo(r *Repository) {
	Repo = r
}

// Home is the about page handler
func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Contacts renders the contacts page and display form
func (rp *Repository) Contacts(writer http.ResponseWriter, request *http.Request) {
	render.RenderTemplate(writer, request, "contacts.page.tmpl", &models.TemplateData{})
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
		helpers.ServerError(writer, err)
		return
	}
	writer.Header().Set("Content-type", "application/json")
	writer.Write(out)
}

// Reservation renders the make a reservation page and display form
func (rp *Repository) Reservation(writer http.ResponseWriter, request *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(writer, request, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation from
func (rp *Repository) PostReservation(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	reservation := models.Reservation{
		FirstName: request.Form.Get("first_name"),
		LastName:  request.Form.Get("last_name"),
		Email:     request.Form.Get("email"),
		Phone:     request.Form.Get("phone"),
	}

	form := forms.New(request.PostForm)
	form.CheckRequiredFields("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(writer, request, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	rp.App.Session.Put(request.Context(), "reservation", reservation)
	http.Redirect(writer, request, "/reservation-summary", http.StatusSeeOther)

}

// ReservationSummary
func (rp *Repository) ReservationSummary(writer http.ResponseWriter, request *http.Request) {
	reservation, ok := rp.App.Session.Pop(request.Context(), "reservation").(models.Reservation)
	if !ok {
		errorTest := "Can not get reservation from session"
		rp.App.ErrorLog.Println(writer, errorTest)
		rp.App.Session.Put(request.Context(), "error", errorTest)
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(writer, request, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
