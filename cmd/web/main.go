package main

import (
	"github.com/alexedwards/scs/v2"
	"github.com/ekateryna-tln/booking/pkg/config"
	"github.com/ekateryna-tln/booking/pkg/hendlers"
	"github.com/ekateryna-tln/booking/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"
var appConfig config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {

	appConfig.UseCache = false
	appConfig.CookieSecure = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.CookieSecure
	appConfig.Session = session

	if appConfig.UseCache {
		templateCache := render.GetTemplateCache()
		appConfig.TemplateCache = templateCache
	}

	render.SetRenderAppConfig(&appConfig)
	hendlers.SetHandlersRepo(hendlers.NewRepo(&appConfig))

	srv := &http.Server{
		Addr: portNumber,
		Handler: routes(&appConfig),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
