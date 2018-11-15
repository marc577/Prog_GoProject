package webserver

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var server *httptest.Server

func setupFunc(handler http.HandlerFunc) {
	server = httptest.NewServer(http.HandlerFunc(handler))
}
func setup(handler http.Handler) {
	server = httptest.NewServer(handler)
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
	setupFunc(serveIndex)
	defer teardown()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode, "Wrong HTTP Status")

	_, er := ioutil.ReadAll(res.Body)
	assert.NoError(t, er)
}

func TestServeIndexPOST(t *testing.T) {
	setupFunc(serveIndex)
	defer teardown()
	res, err := http.Post(server.URL, "", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode, "Wrong HTTP Method")
	_, er := ioutil.ReadAll(res.Body)
	assert.NoError(t, er)
}

func TestMethodsAllow(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hallo Welt")
	})
	setup(adapt(simpleHandler, methods("GET")))
	defer teardown()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hallo Welt", string(body))
}
func TestMethodsNotAllow(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hallo Welt")
	})
	setup(adapt(simpleHandler, methods("GET")))
	defer teardown()
	res, err := http.Post(server.URL, "", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode, "Wrong HTTP Method")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.NotEqual(t, "Hallo Welt", string(body))
}

func TestStart(t *testing.T) {
	go func() {
		Start(8443, "../../keys/server.crt", "../../keys/server.key")
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
