package handlers

import (
	"encoding/json"
	"github.com/ekateryna-tln/booking/internal/config"
	"github.com/ekateryna-tln/booking/internal/driver"
	"github.com/ekateryna-tln/booking/internal/forms"
	"github.com/ekateryna-tln/booking/internal/helpers"
	"github.com/ekateryna-tln/booking/internal/models"
	"github.com/ekateryna-tln/booking/internal/render"
	"github.com/ekateryna-tln/booking/internal/repository"
	"github.com/ekateryna-tln/booking/internal/repository/dbrepo"
	"net/http"
	"time"
)

// Repo the repository used by handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.App
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.App, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// SetHandlersRepo set the repository for the handlers
func SetHandlersRepo(r *Repository) {
	Repo = r
}

// Home is the about page handler
func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Contacts renders the contacts page and display form
func (rp *Repository) Contacts(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "contacts.page.tmpl", &models.TemplateData{})
}

// Generals renders the room page
func (rp *Repository) Generals(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (rp *Repository) Majors(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (rp *Repository) Availability(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	for _, i := range rooms {
		m.App.InfoLog.Println("ROOM:", i.ID, i.RoomName)
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
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
	render.Template(writer, request, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation from
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomID := r.Form.Get("room_id")

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)
	form.CheckRequiredFields("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newResID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restrictions := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newResID,
		RestrictionID: "b1008c53-304a-4226-b16f-ddf8a95c4ed6",
	}

	err = m.DB.InsertRoomRestriction(restrictions)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

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

	render.Template(writer, request, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
