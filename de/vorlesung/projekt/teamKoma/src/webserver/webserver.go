package webserver

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// Authenticator for user autehfication
type Authenticator interface {
	Authenticate(user, password string) bool
}

// The AuthenticatorFunc type is an adapter to allow the use of
// ordinary functions as authenticators.
type AuthenticatorFunc func(user, password string) bool

// Authenticate calls af(user, password).
func (af AuthenticatorFunc) Authenticate(user, password string) bool {
	return af(user, password)
}

func serveIndex(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		t, _ := template.ParseFiles("../../html/index.html")
		t.Execute(w, nil)
	} else {
		http.Error(w, "Method Not Allowed", 405)
		w.WriteHeader(405)
	}

}

// e.g. http.HandleFunc("/health-check", HealthCheckHandler)
// func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
// 	// A very simple health check.
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")

// 	// In the future we could report back on the status of our DB, or our cache
// 	// (e.g. Redis) by performing a simple PING, and include them in the response.
// 	io.WriteString(w, `{"alive": true}`)
// }

// func serveLogin(w http.ResponseWriter, req *http.Request) {
// 	fmt.Println("method:", req.Method) //get request method
// 	if req.Method == "GET" {
// 		t, _ := template.ParseFiles("../../html/login.html")
// 		t.Execute(w, nil)
// 	} else {
// 		req.ParseForm()
// 		// logic part of log in
// 		fmt.Println("username:", req.Form["username"])
// 		fmt.Println("password:", req.Form["password"])
// 	}
// }

func Start(port int, serverCertPath string, serverKeyPath string) error {

	//http Route Handles
	http.HandleFunc("/", serveIndex)
	// http.HandleFunc("/login", serveLogin)

	portString := strings.Join([]string{":", strconv.Itoa(port)}, "")

	httpErr := http.ListenAndServeTLS(portString, serverCertPath, serverKeyPath, nil)

	return httpErr
}

// func basicAuthWrapper(authenticator Authenticator, handler http.HandlerFunc) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user, pswd, ok := r.BasicAuth()
// 		if ok && authenticator.Authenticate(user, pswd) {
// 			handler(w, r)
// 		} else {
// 			w.Header().Set("WWW-Authenticate", "Basic realm=\"KOMA Ticket System\"")
// 			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 		}
// 	})
// }

// func cookieAuthWrapper(authenticator Authenticator, handler http.HandlerFunc) http.Handler {
// 	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	// 	cookie, err := r.Cookie("cookie")
// 	// 	user := cookie.Value
// 	// 	if err == nil && authenticator.Authenticate(user, pswd) {
// 	// 		handler(w, r)
// 	// 	} else {
// 	// 		w.Header().Set("WWW-Authenticate", "Basic realm=\"My Simple Server\"")
// 	// 		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 	// 	}
// 	// })
// }
