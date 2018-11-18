package webserver

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

func (c contextKey) String() string {
	return "webserver_" + string(c)
}

var (
	contextKeyUser = contextKey("user")
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
func mustParams(params ...string) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			for _, param := range params {
				if len(q.Get(param)) == 0 {
					http.Error(w, "missing "+param, http.StatusBadRequest)
					return // exit early
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

func basicAuthWrapper(authenticator Authenticator) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pswd, ok := r.BasicAuth()
			if ok && authenticator.Authenticate(user, pswd) {
				ctx := context.WithValue(r.Context(), contextKeyUser, user)
				h.ServeHTTP(w, r.WithContext(ctx))
			} else {
				w.Header().Set("WWW-Authenticate", "Basic realm=\"KOMA Ticket System\"")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			}
		})
	}
}

func serveTemplate(t *template.Template, name string, data interface{}) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := t.ExecuteTemplate(w, name, data)
			if err != nil {
				h.ServeHTTP(w, r)
			} else {
				//panic
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

var htmlRoot string

func serveIndex(w http.ResponseWriter, req *http.Request) {
	tPath := strings.Join([]string{htmlRoot, "index.html"}, "/")
	t, _ := template.ParseFiles(tPath)
	t.Execute(w, nil)
}
func serveDashAll(w http.ResponseWriter, req *http.Request) {
	tPath := strings.Join([]string{htmlRoot, "dashboardAll.html"}, "/")
	t, _ := template.ParseFiles(tPath)
	t.Execute(w, nil)
}
func serveDashUn(w http.ResponseWriter, req *http.Request) {
	tPath := strings.Join([]string{htmlRoot, "dashboardUnassigned.html"}, "/")
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

	// static files
	staticFilePath := htmlRoot + "/" + "assets"
	fs := http.FileServer(http.Dir(staticFilePath))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// templates
	tmpls := make(map[string]*template.Template)
	tmpls["index"] = template.Must(template.ParseFiles(rootPath+"/newTicket.tmpl.html", rootPath+"/layout.tmpl.html"))
	tmpls["dash"] = template.Must(template.ParseFiles(rootPath+"/dashAll.tmpl.html", rootPath+"/layout.tmpl.html"))

	// templates, err := template.ParseFiles(allFiles...)
	// if err != nil {
	// 	//panic
	// }
	// tree := templates.Tree
	// name := templates.DefinedTemplates()
	// layoutTmpl := templates.Lookup("layout")
	// fmt.Print(tree, name)
	// newTicketTmpl := templates.Lookup("newticket.tmpl.html")
	// layoutTmpl := templates.Lookup("layout.tmpl.html")
	// layoutTmpl := template.Must(templates.Lookup("layout.tmpl.html"))
	// layoutTmpl.ExecuteTemplate(os.Stdout, "layout", nil)
	// fmt.Println()

	//newTicketTemplate := templates.Lookup("newTicket.html")
	//dashAll := templates.Lookup("dashAll.html")

	// auth := func(user, pswd string) bool {
	// 	fmt.Println(user, ":", pswd)
	// 	return true
	// }

	http.Handle("/dash", adapt(nil, serveTemplate(tmpls["dash"], "layout", nil)))
	// http.Handle("/dashu", adapt(nil, serveTemplate(tmpls["dash"], "layout", nil)))
	// http.Handle("/dasho", adapt(nil, serveTemplate(tmpls["dash"], "layout", nil)))
	// http.Handle("/dash", adapt(http.HandlerFunc(serveDashAll), basicAuthWrapper(AuthenticatorFunc(auth))))
	http.Handle("/", adapt(nil, serveTemplate(tmpls["index"], "layout", nil)))

	portString := ":" + strconv.Itoa(port)

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
