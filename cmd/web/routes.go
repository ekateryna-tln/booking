package main

import (
	"github.com/ekateryna-tln/booking/internal/config"
	"github.com/ekateryna-tln/booking/internal/hendlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func routes(appConfig *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(LoadSession)

	mux.Get("/", hendlers.Repo.Home)
	mux.Get("/about", hendlers.Repo.About)
	mux.Get("/contacts", hendlers.Repo.Contacts)
	mux.Get("/generals-quarters", hendlers.Repo.Generals)
	mux.Get("/majors-suites", hendlers.Repo.Majors)
	mux.Get("/search-availability", hendlers.Repo.Availability)
	mux.Post("/search-availability", hendlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", hendlers.Repo.AvailabilityJSON)
	mux.Get("/make-reservation", hendlers.Repo.Reservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
