package webserver

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type contextKey string

func (c contextKey) String() string {
	return "webserver_" + string(c)
}

var (
	contextKeyUser = contextKey.String("user")
)

type adapter func(http.HandlerFunc) http.HandlerFunc

func methodsWrapper(methods ...string) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			method := r.Method
			for _, me := range methods {
				if me == method {
					if h != nil {
						h.ServeHTTP(w, r)
					}
					return
				}
			}
			http.Error(w, "Method Not Allowd", 405)
			w.WriteHeader(405)
		})
	}
}
func mustParamsWrapper(params ...string) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			for _, param := range params {
				if len(r.Form.Get(param)) == 0 {
					http.Error(w, "missing "+param, http.StatusBadRequest)
					return
				}
			}
			if h != nil {
				h.ServeHTTP(w, r)
			}
		})
	}
}

func basicAuthWrapper(authenticator Authenticator) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pswd, ok := r.BasicAuth()
			if ok && authenticator.Authenticate(user, pswd) {
				ctx := context.WithValue(r.Context(), contextKey("user"), user)
				if h != nil {
					h.ServeHTTP(w, r.WithContext(ctx))
				}
			} else {
				w.Header().Set("WWW-Authenticate", "Basic realm=\"KOMA Ticket System\"")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			}
		})
	}
}

func serveTemplateWrapper(t *template.Template, name string, data interface{}) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := t.ExecuteTemplate(w, name, data)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			if h != nil {
				h.ServeHTTP(w, r)
			}
		})
	}
}

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

func Start(port int, serverCertPath string, serverKeyPath string, rootPath string) error {

	htmlRoot := rootPath

	// static files
	staticFilePath := htmlRoot + "/" + "assets"
	fs := http.FileServer(http.Dir(staticFilePath))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// templates
	tmpls := make(map[string]*template.Template)
	tmpls["index"] = template.Must(template.ParseFiles(rootPath+"/new.tmpl.html", rootPath+"/index.tmpl.html"))
	tmpls["admin"] = template.Must(template.ParseFiles(rootPath+"/dashboard.tmpl.html", rootPath+"/index.tmpl.html"))
	tmpls["new"] = template.Must(template.ParseFiles(rootPath+"/new.tmpl.html", rootPath+"/index.tmpl.html"))

	auth := AuthenticatorFunc(func(user, pswd string) bool {
		fmt.Println(user, ":", pswd)
		return true
	})

	// frontend
	http.Handle("/admin", adapt(nil, serveTemplateWrapper(tmpls["admin"], "layout", nil)))
	http.Handle("/assigned", adapt(nil, serveTemplateWrapper(tmpls["admin"], "layout", nil)))
	http.Handle("/all", adapt(nil, serveTemplateWrapper(tmpls["admin"], "layout", nil)))
	http.Handle("/new", adapt(nil, serveTemplateWrapper(tmpls["new"], "layout", nil)))
	http.Handle("/", adapt(nil, serveTemplateWrapper(tmpls["index"], "layout", nil)))

	// rest-api
	// insert ticket via mail
	http.Handle("/api/new", adapt(nil, mustParamsWrapper("POST"), methodsWrapper("POST"), basicAuthWrapper(auth)))
	// mail sending
	http.Handle("/api/mail", adapt(nil, mustParamsWrapper("POST"), methodsWrapper("GET"), basicAuthWrapper(auth)))
	http.Handle("/api/mail", adapt(nil, mustParamsWrapper("POST"), methodsWrapper("POST"), basicAuthWrapper(auth)))

	portString := ":" + strconv.Itoa(port)
	httpErr := http.ListenAndServeTLS(portString, serverCertPath, serverKeyPath, nil)

	return httpErr
}
