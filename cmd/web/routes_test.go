package main

import (
	"github.com/ekateryna-tln/booking/internal/config"
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRouts(t *testing.T) {
	var app config.App
	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing; Test passed
	default:
		t.Errorf("type is not *chi.Mux %T", v)
	}
}
