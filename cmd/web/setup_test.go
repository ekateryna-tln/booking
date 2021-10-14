package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type plugHandler struct{}

func (mh *plugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
