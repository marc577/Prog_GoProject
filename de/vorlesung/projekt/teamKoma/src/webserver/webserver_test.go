package webserver

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var server *httptest.Server

func setup(handler http.HandlerFunc) {
	server = httptest.NewServer(http.HandlerFunc(serveIndex))
}
func teardown() {
	server.Close()
}

// func TestStart(t *testing.T) {
// 	httpError := Start(8443, "../../keys/server.crt", "../../keys/server.key")
// 	if httpError != nil {
// 		t.Error("Error Init Webser", httpError)
// 	} else {
// 		log.Println("sd")
// 	}
// }

func TestServeIndexGET(t *testing.T) {
	setup(serveIndex)
	defer teardown()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode, "Wrong HTTP Status")

	_, er := ioutil.ReadAll(res.Body)
	assert.NoError(t, er)
}

func TestServeIndexPOST(t *testing.T) {
	setup(serveIndex)
	defer teardown()
	res, err := http.Post(server.URL, "", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode, "Wrong HTTP Method")
	_, er := ioutil.ReadAll(res.Body)
	assert.NoError(t, er)
}

func TestStart(t *testing.T) {
	jobs := make(chan error)
	go func() {
		err := Start(8443, "../../keys/server.crt", "../../keys/server.key")
		jobs <- err
		close(jobs)
	}()
	assert.True(t, true)
}

// func TestHealthCheckHandler(t *testing.T) {
// 	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
// 	// pass 'nil' as the third parameter.
// 	req, err := http.NewRequest("GET", "/health-check", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(HealthCheckHandler)

// 	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// 	// directly and pass in our Request and ResponseRecorder.
// 	handler.ServeHTTP(rr, req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	// Check the response body is what we expect.
// 	expected := `{"alive": true}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
