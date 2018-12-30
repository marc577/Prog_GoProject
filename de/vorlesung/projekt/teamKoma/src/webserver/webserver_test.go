//Matrikelnummern:
//9188103
//1798794
//4717960
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
	"net/url"
	"storagehandler"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var server *httptest.Server

func setup(handler http.Handler) storagehandler.StorageWrapper {
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
		"getUser":    func() string { return "" },
		"getHoliday": func(name string) bool { return false },
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

// from https://www.dotnetperls.com/between-before-after-go
func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}
func TestStart(t *testing.T) {

	urlsGET := []string{"/", "/open", "/assigned", "/all"}

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
	// /new
	client := &http.Client{Transport: transCfg}
	form := url.Values{}
	form.Add("lName", "a")
	form.Add("fName", "b")
	form.Add("email", "a@b")
	form.Add("subject", "Test")
	form.Add("description", "Das ist ein Test")
	req, err := http.NewRequest("POST", host+"/new", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	s := string(body)
	tID := between(s, "<h3>", "</h3>")
	assert.NoError(t, err)
	assert.NotEqual(t, "", tID)
	assert.Equal(t, http.StatusOK, res.StatusCode, tID)

	// /edit
	form = url.Values{}
	form.Add("state", "1")
	form.Add("processor", "Werner")
	req, err = http.NewRequest("POST", host+"/edit?ticket="+tID, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("Werner", "password")
	res, err = client.Do(req)
	assert.NoError(t, err)
	body, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode, string(body))

	// /edit/add
	form = url.Values{}
	form.Add("type", "Inform")
	form.Add("description", "Test Item")
	form.Add("email", "a@k")
	req, err = http.NewRequest("POST", host+"/edit/add?ticket="+tID, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("Werner", "password")
	res, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// /edit/free
	req, err = http.NewRequest("GET", host+"/edit/free?ticket="+tID, nil)
	req.SetBasicAuth("Werner", "password")
	res, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// /user/holiday
	req, err = http.NewRequest("POST", host+"/user/holiday", bytes.NewReader(body))
	req.SetBasicAuth("Werner", "password")
	res, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// /edit/assign
	req, err = http.NewRequest("GET", host+"/assign?ticket="+tID+"&user=Werner", nil)
	req.SetBasicAuth("Werner", "password")
	res, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// /edit/mail
	req, err = http.NewRequest("GET", host+"/api/mail", nil)
	req.SetBasicAuth("Werner", "password")
	res, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	body, err = ioutil.ReadAll(res.Body)
	assert.NoError(t, err, s)

	// /edit/mail
	req, err = http.NewRequest("POST", host+"/api/mail", bytes.NewReader(body))
	req.SetBasicAuth("Werner", "password")
	res, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// /api/new
	// /edit/combine

}
