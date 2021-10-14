package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/ekateryna-tln/booking/internal/config"
	"github.com/ekateryna-tln/booking/internal/handlers"
	"github.com/ekateryna-tln/booking/internal/models"
	"github.com/ekateryna-tln/booking/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.App
var session *scs.SessionManager

// main is the main application function
func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	gob.Register(models.Reservation{})

	app.UseCache = false
	app.CookieSecure = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.CookieSecure
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("create template cache error:", err)
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false

	render.SetRenderApp(&app)
	handlers.SetHandlersRepo(handlers.NewRepo(&app))

	return nil
}
