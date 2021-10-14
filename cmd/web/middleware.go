package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// LoadSession loads and save the session on every request
func LoadSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
