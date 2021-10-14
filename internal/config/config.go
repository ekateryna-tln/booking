package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

// App holds the application config and data
type App struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	CookieSecure  bool
	Session       *scs.SessionManager
}
