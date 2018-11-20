package webserver

import (
	"context"
	"html/template"
	"net/http"
	"storagehandler"
	"strconv"
)

type contextKey string

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

func dataWrapperAll() adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := storagehandler.GetTickets()
			ctx := context.WithValue(r.Context(), contextKey("data"), data)
			if h != nil {
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}
func dataWrapperOpen() adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := storagehandler.GetOpenTickets()
			ctx := context.WithValue(r.Context(), contextKey("data"), data)
			if h != nil {
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}
func dataWrapperAssigned() adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxVal := r.Context().Value(contextKey("user"))
			if ctxVal != nil {
				user := ctxVal.(string)
				data := storagehandler.GetNotClosedTicketsByProcessor(user)
				ctx := context.WithValue(r.Context(), contextKey("data"), data)
				if h != nil {
					h.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		})
	}
}
func serveTemplateWrapper(t *template.Template, name string, data interface{}) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if data == nil {
				ctxKey := contextKey("data")
				ctxVal := r.Context().Value(ctxKey)
				data = ctxVal
			}
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

// Start initializes the webserver with the required
// parameters, registers the urls and sets the Authenticator
// function to the VerifyUser function
// from the storagehandler packet
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

	auth := AuthenticatorFunc(storagehandler.VerifyUser)

	// frontend
	http.Handle("/admin", adapt(nil, serveTemplateWrapper(tmpls["admin"], "layout", nil), dataWrapperOpen(), basicAuthWrapper(auth), methodsWrapper("GET")))
	http.Handle("/assigned", adapt(nil, serveTemplateWrapper(tmpls["admin"], "layout", nil), dataWrapperAssigned(), basicAuthWrapper(auth), methodsWrapper("GET")))
	http.Handle("/all", adapt(nil, serveTemplateWrapper(tmpls["admin"], "layout", nil), dataWrapperAll(), basicAuthWrapper(auth), methodsWrapper("GET")))
	http.Handle("/new", adapt(nil, serveTemplateWrapper(tmpls["new"], "layout", nil)))
	http.Handle("/", adapt(nil, serveTemplateWrapper(tmpls["index"], "layout", nil)))

	// rest-api
	// insert ticket via mail
	http.Handle("/api/new", adapt(nil, mustParamsWrapper("POST"), basicAuthWrapper(auth), methodsWrapper("POST")))
	// mail sending
	http.Handle("/api/mail", adapt(nil, mustParamsWrapper("POST"), basicAuthWrapper(auth), methodsWrapper("GET", "POST")))

	//http.Handle("/api/mail", adapt(nil, mustParamsWrapper("POST"), methodsWrapper("POST"), basicAuthWrapper(auth)))

	portString := ":" + strconv.Itoa(port)
	httpErr := http.ListenAndServeTLS(portString, serverCertPath, serverKeyPath, nil)

	return httpErr
}
