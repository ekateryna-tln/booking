package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	plugH := plugHandler{}
	h := NoSurf(&plugH)
	switch v := h.(type) {
	case http.Handler:
		//do nothing; Test passed
	default:
		t.Errorf("type is not http.Handler %T", v)
	}
}
func TestLoadSession(t *testing.T) {
	plugH := plugHandler{}
	h := LoadSession(&plugH)
	switch v := h.(type) {
	case http.Handler:
		//do nothing; Test passed
	default:
		t.Errorf("type is not http.Handler %T", v)
	}
}
