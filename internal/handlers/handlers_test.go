package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ekateryna-tln/booking/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "get", http.StatusOK},
	{"about", "/about", "get", http.StatusOK},
	{"gq", "/generals-quarters", "get", http.StatusOK},
	{"ms", "/majors-suite", "get", http.StatusOK},
	{"sa", "/search-availability", "get", http.StatusOK},
	{"contacts", "/contacts", "get", http.StatusOK},
	{"mr", "/make-reservation", "get", http.StatusOK},
	{"rs", "/reservation-summary", "get", http.StatusOK},
	//{"post_sa", "/search-availability", "post", []postData{}, http.StatusOK},
	//{"post_sa", "/search-availability", "post", []postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-01-10"},
	//}, http.StatusOK},
	//{"post_sa_json", "/search-availability-json", "post", []postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-01-10"},
	//}, http.StatusOK},
	//{"post_mr", "/make-reservation", "post", []postData{
	//	{key: "first_name", value: "test_first_name"},
	//	{key: "last_name", value: "test_last_name"},
	//	{key: "email", value: "test@test.test"},
	//	{key: "phone", value: "12342341234"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range theTests {
		response, err := testServer.Client().Get(testServer.URL + test.url)
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

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: "f1ef9c13-1aa8-4245-b128-3d479d1b87a2",
		Room: models.Room{
			ID:       "f1ef9c13-1aa8-4245-b128-3d479d1b87a2",
			RoomName: "General's Quaters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)

	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	reservation.RoomID = ""
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func TestRepository_PostReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID:    "f1ef9c13-1aa8-4245-b128-3d479d1b87a2",
		StartDate: time.Date(3000, 1, 1, 0, 0, 0, 0, time.Local),
		EndDate:   time.Date(3000, 1, 10, 0, 0, 0, 0, time.Local),
		Room: models.Room{
			ID:       "f1ef9c13-1aa8-4245-b128-3d479d1b87a2",
			RoomName: "General's Quaters",
		},
	}

	postedDate := url.Values{}
	postedDate.Add("first_name", "TestFirstName")
	postedDate.Add("last_name", "TestLastName")
	postedDate.Add("email", "test@test.test")
	postedDate.Add("phone", "123456789")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx := getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for success case: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure to insert reservation into database
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx = getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	reservation.RoomID = ""
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for success case: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for failure to insert room restriction into database
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx = getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	reservation.RoomID = "fail-room-restriction-uuid"
	reservation.Room.ID = "fail-room-restriction-uuid"
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when reservation insert error happened: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case where reservation is not in session
	req, _ = http.NewRequest("GET", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for empty session data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing form body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid data
	postedDate.Set("first_name", "a")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedDate.Encode()))
	ctx = getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// first case - rooms are available
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=f1ef9c13-1aa8-4245-b128-3d479d1b87a2")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	var jsonResponse jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &jsonResponse)
	if err != nil {
		t.Error("failed to parse json")
	}
	if jsonResponse.OK != true {
		t.Error("AvailabilityJSON returns room not available when it is")
	}

	// second case - rooms are not available
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=room_not_available")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &jsonResponse)
	if err != nil {
		t.Error("failed to parse json")
	}
	if jsonResponse.OK != false {
		t.Error("AvailabilityJSON returns room available when it is not")
	}

	// third case - database error during search availability
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=fail-search-availability")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &jsonResponse)
	if err != nil {
		t.Error("failed to parse json")
	}
	if jsonResponse.Message != "Error connecting to database" {
		t.Error("AvailabilityJSON handler failed when search availability error happened")
	}

	// fourth case - start date parse error
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=fail-search-availability")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &jsonResponse)
	if err != nil {
		t.Error("failed to parse json")
	}
	if jsonResponse.Message != "StartDate parse error" {
		t.Error("AvailabilityJSON handler failed when start date parse error happened")
	}

	// fifth case - end date parse error
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=fail-search-availability")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &jsonResponse)
	if err != nil {
		t.Error("failed to parse json")
	}
	if jsonResponse.Message != "EndDate parse error" {
		t.Error("AvailabilityJSON handler failed when end date parse error happened")
	}
}
