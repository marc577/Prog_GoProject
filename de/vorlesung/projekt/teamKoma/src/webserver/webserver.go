package webserver

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Authenticator interface {
	Authenticate(user, password string) bool
}

type AuthenticatorFunc func(user, password string) bool

func (af AuthenticatorFunc) Authenticate(user, password string) bool {
	return af(user, password)
}

func serveIndex(w http.ResponseWriter, req *http.Request) {
	AuthenticatorFunc.Authenticate(func(name, pwd string) bool {
		return true
	}, "sd", "sd")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

// e.g. http.HandleFunc("/health-check", HealthCheckHandler)
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

func serveLogin(w http.ResponseWriter, req *http.Request) {
	fmt.Println("method:", req.Method) //get request method
	if req.Method == "GET" {
		t, _ := template.ParseFiles("../../html/login.html")
		t.Execute(w, nil)
	} else {
		req.ParseForm()
		// logic part of log in
		fmt.Println("username:", req.Form["username"])
		fmt.Println("password:", req.Form["password"])
	}
}

func Init(port int, serverCertPath string, serverKeyPath string) error {

	//http Route Handles
	http.HandleFunc("/", serveIndexHandler)
	http.HandleFunc("/login", serveLogin)
	http.HandleFunc("/login2", secureHandlerWith(,serveLogin))

	portString := strings.Join([]string{":", strconv.Itoa(port)}, "")

	httpErr := http.ListenAndServeTLS(portString, serverCertPath, serverKeyPath, nil)

	if httpErr != nil {
		log.Println("ListenAndServe: ", httpErr)
	}
	log.Printf("ListenAndServe on Port %d \n", port)

	return httpErr
}

func secureHandlerWith(authenticator Authenticator, handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pswd, ok := r.BasicAuth()
		if ok && authenticator.Authenticate(user, pswd) {
			handler(w, r)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"My Simple Server\"")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	})
}
