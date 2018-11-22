package webserver

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"storagehandler"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var server *httptest.Server

func setupFunc(handler http.HandlerFunc) {
	storagehandler.Init("../../storage/users.json", "../../storage/tickets")
	server = httptest.NewServer(http.HandlerFunc(handler))
}
func setup(handler http.Handler) {
	storagehandler.Init("../../storage/users.json", "../../storage/tickets")
	server = httptest.NewServer(handler)
}
func setupSimple(handler http.HandlerFunc) {
	server = httptest.NewServer(handler)
}
func teardown() {
	server.Close()
}

func TestServeTemplate(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	})
	rootPath := "../../html"
	tmpl := template.Must(template.ParseFiles(rootPath+"/new.tmpl.html", rootPath+"/index.tmpl.html"))
	setup(adapt(simpleHandler, serveTemplateWrapper(tmpl, "layout", nil), dataWrapperAll()))
	defer teardown()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode, "Wrong HTTP Status")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}
func TestServeTemplateFalse(t *testing.T) {
	rootPath := "../../html"
	tmpl := template.Must(template.ParseFiles(rootPath+"/new.tmpl.html", rootPath+"/index.tmpl.html"))
	setup(adapt(nil, serveTemplateWrapper(tmpl, "layout2", nil)))
	defer teardown()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode, "Wrong HTTP Status")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Empty(t, body)
}

func TestMethodsAllow(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hallo Welt")
	})
	setup(adapt(simpleHandler, methodsWrapper("GET")))
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
	setup(adapt(simpleHandler, methodsWrapper("GET")))
	defer teardown()
	res, err := http.Post(server.URL, "", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode, "Wrong HTTP Method")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.NotEqual(t, "Hallo Welt", string(body))
}
func TestMustParamsOK(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		responseString := strings.Join([]string{"Hallo", req.URL.Query().Get("fname"), req.Form.Get("lname")}, " ")
		io.WriteString(w, responseString)
	})
	setup(adapt(simpleHandler, mustParamsWrapper("fname", "lname")))
	defer teardown()
	url := strings.Join([]string{server.URL, "fname=Werner"}, "?")
	res, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString("lname=Brenzel"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hallo Werner Brenzel", string(body))
}
func TestMustParamsNotOK(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		responseString := strings.Join([]string{"Hallo", req.URL.Query().Get("name")}, " ")
		io.WriteString(w, responseString)
	})
	setup(adapt(simpleHandler, mustParamsWrapper("name")))
	defer teardown()
	url := strings.Join([]string{server.URL, "greet=Werner"}, "?")
	res, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.NotEqual(t, "Hallo Werner", string(body))
}

func TestBasicAuthWrapperWithoutPW(t *testing.T) {
	setup(adapt(nil, basicAuthWrapper(nil)))
	defer teardown()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "wrong status")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized)+"\n", string(body), "wrong message")
}

func TestBasicAuthWrapperWithOKPW(t *testing.T) {
	var receivedName, receivedPW string
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctxKey := contextKey("user")
		ctxVal := req.Context().Value(ctxKey)
		assert.Equal(t, "<username>", ctxVal.(string), "Context not set")
		io.WriteString(w, "Hello client\n")
	})
	auth := AuthenticatorFunc(func(n string, p string) bool {
		receivedName = n
		receivedPW = p
		return true
	})
	setup(adapt(simpleHandler, basicAuthWrapper(auth)))
	client := &http.Client{}
	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)
	req.SetBasicAuth("<username>", "<password>")
	res, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode, "wrong status code")
	assert.Equal(t, "<username>", receivedName, "wrong username")
	assert.Equal(t, "<password>", receivedPW, "wrong password")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello client\n", string(body), "wrong message")
}

func TestBasicAuthWrapperWithNotOKPW(t *testing.T) {
	var receivedName, receivedPW string
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello client\n")
	})
	auth := AuthenticatorFunc(func(n string, p string) bool {
		receivedName = n
		receivedPW = p
		return false
	})
	setup(adapt(simpleHandler, basicAuthWrapper(auth)))
	client := &http.Client{}
	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)
	req.SetBasicAuth("<username>", "<password>")
	res, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "wrong status code")
	assert.Equal(t, "<username>", receivedName, "wrong username")
	assert.Equal(t, "<password>", receivedPW, "wrong password")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized)+"\n", string(body), "wrong message")
}

// func TestRedirectWrapper(t *testing.T) {
// 	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		io.WriteString(w, "Hello client\n")
// 	})
// 	server.Config.
// 	setupSimple(adapt(simpleHandler, redirectWrapper("edit")))
// 	defer teardown()
// 	res, err := http.Get(server.URL)
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusPermanentRedirect, res.StatusCode, "wrong status")
// 	body, err := ioutil.ReadAll(res.Body)
// 	assert.NoError(t, err)
// 	assert.NotEqual(t, "Hello client\n", string(body), "not redirectec")
// }

func TestStart(t *testing.T) {
	//c := make(chan error, 1)
	go func() {
		Start(8443, "../../keys/server.crt", "../../keys/server.key", "../../html")
		//close(c)
	}()
	//time.Sleep(2 * time.Second)
	//err := <-c
	//assert.NoError(t, err)
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
