package webserver

import (
	"bytes"
	"context"
	"crypto/tls"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"storagehandler"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var server *httptest.Server

func setup(handler http.Handler) *storagehandler.StorageHandler {
	server = httptest.NewServer(handler)
	return storagehandler.New("../../storage/users.json", "../../storage/tickets")
}
func setupSimple(handler http.Handler) {
	server = httptest.NewServer(handler)
}
func teardown() {
	server.Close()
}

func TestVerifyMail(t *testing.T) {
	assert.True(t, verifyEMail("ale@kale"))
	assert.True(t, verifyEMail("ale@kale.de"))
	assert.False(t, verifyEMail("ale@k!.--slale.de"))
	assert.False(t, verifyEMail("ale.--slale.de"))
	assert.False(t, verifyEMail("HalloWelt"))
}
func TestMethodsAllow(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hallo Welt")
	})
	setupSimple(adapt(simpleHandler, methodsWrapper("GET")))
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
	setupSimple(adapt(simpleHandler, methodsWrapper("GET")))
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
	setupSimple(adapt(simpleHandler, mustParamsWrapper("fname", "lname")))
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
	setupSimple(adapt(simpleHandler, mustParamsWrapper("name")))
	defer teardown()
	url := strings.Join([]string{server.URL, "greet=Werner"}, "?")
	res, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.NotEqual(t, "Hallo Werner", string(body))
}
func TestSaveParamsOK(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		n := req.Context().Value(contextKey("name")).(string)
		responseString := strings.Join([]string{"Hallo", n}, " ")
		io.WriteString(w, responseString)
	})
	setupSimple(adapt(simpleHandler, saveParamsWrapper("name")))
	defer teardown()
	url := strings.Join([]string{server.URL, "name=Werner"}, "?")
	res, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hallo Werner", string(body))
}
func TestSaveParamsNotOK(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		n := req.Context().Value(contextKey("name"))
		if n == nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		responseString := strings.Join([]string{"Hallo", "Werner"}, " ")
		io.WriteString(w, responseString)
	})
	setupSimple(adapt(simpleHandler, saveParamsWrapper("name2")))
	defer teardown()
	url := strings.Join([]string{server.URL, "name=Werner"}, "?")
	res, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.NotEqual(t, "Hallo Werner", string(body))
}
func TestRedirectWrapper(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hallo")
	})
	setupSimple(adapt(simpleHandler, redirectWrapper("")))
	defer teardown()
	res, err := http.Get(server.URL)
	assert.Error(t, err)
	assert.Nil(t, res)
}
func TestFunctionCTXWrapper(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctxval := req.Context().Value(contextKey("data"))
		assert.NotNil(t, ctxval)
		greet := "Hallo " + ctxval.(string)
		io.WriteString(w, greet)
	})
	setupSimple(adapt(simpleHandler, functionCtxWrapper(func(w http.ResponseWriter, r *http.Request) context.Context {
		return context.WithValue(r.Context(), contextKey("data"), "Werner")
	})))
	defer teardown()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hallo Werner", string(body))
}
func TestFunctionCTXWrapperNil(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctxval := req.Context().Value(contextKey("data"))
		assert.Nil(t, ctxval)
	})
	setupSimple(adapt(simpleHandler, functionCtxWrapper(func(w http.ResponseWriter, r *http.Request) context.Context {
		return nil
	})))
	defer teardown()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
func TestServeTemplate(t *testing.T) {
	simpleHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	})
	defaultOpenT := template.New("").Funcs(map[string]interface{}{
		"getUser": func() string { return "" },
	})
	rootPath := "../../html"
	tmpl := template.Must(defaultOpenT.ParseFiles(rootPath+"/orow.tmpl.html", rootPath+"/dashboard.tmpl.html", rootPath+"/index.tmpl.html"))
	setupSimple(adapt(simpleHandler, serveTemplateWrapper(tmpl, "layout", nil), functionCtxWrapper(func(w http.ResponseWriter, r *http.Request) context.Context {
		return context.WithValue(r.Context(), contextKey("user"), "Werner")
	})))
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
	setupSimple(adapt(nil, serveTemplateWrapper(tmpl, "layout2", nil)))
	defer teardown()
	res, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode, "Wrong HTTP Status")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Empty(t, body)
}

func TestBasicAuthWrapperWithoutPW(t *testing.T) {
	setupSimple(adapt(nil, basicAuthWrapper(nil)))
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
	setupSimple(adapt(simpleHandler, basicAuthWrapper(auth)))
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
	setupSimple(adapt(simpleHandler, basicAuthWrapper(auth)))
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

func TestStart(t *testing.T) {

	urlsGET := []string{"/open", "/assigned", "/all", "/api/mail"}

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	host := "https://localhost:8443"
	st := storagehandler.New("../../storage/users.json", "../../storage/tickets/")
	go func() {
		serr := Start(8443, "../../keys/server.crt", "../../keys/server.key", "../../html", st)
		assert.NoError(t, serr)
	}()
	time.Sleep(2 * time.Second)
	for _, url := range urlsGET {
		client := &http.Client{Transport: transCfg}
		req, err := http.NewRequest("GET", host+url, nil)
		req.SetBasicAuth("Werner", "password")
		res, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode, url)
	}

}
