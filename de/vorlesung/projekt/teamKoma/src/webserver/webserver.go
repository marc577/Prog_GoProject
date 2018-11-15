package webserver

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type adapter func(http.HandlerFunc) http.HandlerFunc

func methods(methods ...string) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			method := r.Method
			for _, me := range methods {
				if me == method {
					h.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Method Not Allowd", 405)
			w.WriteHeader(405)
		})
	}
}
func logger() adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RequestURI)
			h.ServeHTTP(w, r)
		})
	}
}
func MustParams() adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RequestURI)
			h.ServeHTTP(w, r)
		})
	}
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

// Adapts several http handlers
// Idea from https://www.youtube.com/watch?v=tIm8UkSf6RA&t=537s
func adapt(h http.HandlerFunc, adapters ...adapter) http.HandlerFunc {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

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

var htmlRoot string

func serveIndex(w http.ResponseWriter, req *http.Request) {
	tPath := strings.Join([]string{htmlRoot, "index.html"}, "/")
	t, _ := template.ParseFiles(tPath)
	t.Execute(w, nil)
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

func Start(port int, serverCertPath string, serverKeyPath string, rootPath string) error {

	htmlRoot = rootPath

	staticFilePath := strings.Join([]string{htmlRoot, "assets"}, "/")

	//http Route Handles
	fs := http.FileServer(http.Dir(staticFilePath))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.Handle("/", adapt(http.HandlerFunc(serveIndex), logger()))

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
