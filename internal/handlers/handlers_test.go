package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "get", []postData{}, http.StatusOK},
	{"about", "/about", "get", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "get", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "get", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "get", []postData{}, http.StatusOK},
	{"contacts", "/contacts", "get", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "get", []postData{}, http.StatusOK},
	{"rs", "/reservation-summary", "get", []postData{}, http.StatusOK},
	{"post_sa", "/search-availability", "post", []postData{}, http.StatusOK},
	{"post_sa", "/search-availability", "post", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-10"},
	}, http.StatusOK},
	{"post_sa_json", "/search-availability-json", "post", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-10"},
	}, http.StatusOK},
	{"post_mr", "/make-reservation", "post", []postData{
		{key: "first_name", value: "test_first_name"},
		{key: "last_name", value: "test_last_name"},
		{key: "email", value: "test@test.test"},
		{key: "phone", value: "12342341234"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range theTests {
		if test.method == "get" {
			response, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("for the test %s expected status code is %d, but got %d",
					test.name, test.expectedStatusCode, response.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, data := range test.params {
				values.Add(data.key, data.value)
			}
			response, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("for the test %s expected status code is %d, but got %d",
					test.name, test.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
