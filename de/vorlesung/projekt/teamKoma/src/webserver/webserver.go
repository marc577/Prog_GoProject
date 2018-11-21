package webserver

import (
	"context"
	"html/template"
	"net/http"
	"regexp"
	"storagehandler"
	"strconv"
)

type contextKey string

type adapter func(http.HandlerFunc) http.HandlerFunc

// verfiy if a string is a valid email adress
// from: http://www.golangprograms.com/regular-expression-to-validate-email-address.html
func verifyEMail(mail string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(mail)
}

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

func newTicketWrapper() adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.PostForm == nil {
				r.ParseForm()
			}
			if verifyEMail(r.Form.Get("email")) != true {
				http.Error(w, http.StatusText(http.StatusNotFound)+"|email", http.StatusNotFound)
			} else {
				t := storagehandler.CreateTicket(r.Form.Get("subject"), r.Form.Get("email"), r.Form.Get("description"))
				if h != nil {
					ctx := context.WithValue(r.Context(), contextKey("data"), t)
					h.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		})
	}
}

// func adaptFuncWrapper(f http.HandlerFunc) adapter {
// 	return func(h http.HandlerFunc) http.HandlerFunc {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			if f != nil {
// 				f.ServeHTTP(w, r)
// 			}
// 			if h != nil {
// 				h.ServeHTTP(w, r)
// 			}
// 		})
// 	}
// }

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

type webContext struct {
	Data interface{}
	Path string
	User interface{}
}

func dataWrapperOne() adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//web := webContext{storagehandler.GetTicket("sd"), "one", user}
			tID := r.URL.Query().Get("ticket")
			ctx := context.WithValue(r.Context(), contextKey("data"), storagehandler.GetTicket(tID))
			if h != nil {
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}
func dataWrapperAll() adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// user := ""
			// ctxVal := r.Context().Value(contextKey("user"))
			// if ctxVal != nil {
			// 	user = ctxVal.(string)
			// }
			//web := webContext{storagehandler.GetTicketsPointer(), "all", user}
			ctx := context.WithValue(r.Context(), contextKey("data"), storagehandler.GetTicketsPointer())
			if h != nil {
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}
func dataWrapperOpen() adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//web := webContext{storagehandler.GetOpenTickets(), "open", user}
			ctx := context.WithValue(r.Context(), contextKey("data"), storagehandler.GetOpenTickets())
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
				//web := webContext{storagehandler.GetNotClosedTicketsByProcessor(user), "assigned", user}
				ctx := context.WithValue(r.Context(), contextKey("data"), storagehandler.GetNotClosedTicketsByProcessor(user))
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
			user := ""
			ctxVal := r.Context().Value(contextKey("user"))
			if ctxVal != nil {
				user = ctxVal.(string)
			}
			path := r.URL.Path
			web := webContext{data, path, user}
			err := t.ExecuteTemplate(w, name, web)
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
	tmpls["open"] = template.Must(template.ParseFiles(rootPath+"/orow.tmpl.html", rootPath+"/dashboard.tmpl.html", rootPath+"/index.tmpl.html"))
	tmpls["admin"] = template.Must(template.ParseFiles(rootPath+"/row.tmpl.html", rootPath+"/dashboard.tmpl.html", rootPath+"/index.tmpl.html"))
	tmpls["added"] = template.Must(template.ParseFiles(rootPath+"/added.tmpl.html", rootPath+"/index.tmpl.html"))
	tmpls["edit"] = template.Must(template.ParseFiles(rootPath+"/ticket.tmpl.html", rootPath+"/index.tmpl.html"))

	auth := AuthenticatorFunc(storagehandler.VerifyUser)

	// frontend
	http.Handle("/open", adapt(nil, serveTemplateWrapper(tmpls["open"], "layout", nil), dataWrapperOpen(), basicAuthWrapper(auth), methodsWrapper("GET")))
	http.Handle("/assigned", adapt(nil, serveTemplateWrapper(tmpls["admin"], "layout", nil), dataWrapperAssigned(), basicAuthWrapper(auth), methodsWrapper("GET")))
	http.Handle("/all", adapt(nil, serveTemplateWrapper(tmpls["admin"], "layout", nil), dataWrapperAll(), basicAuthWrapper(auth), methodsWrapper("GET")))
	http.Handle("/new", adapt(nil, serveTemplateWrapper(tmpls["added"], "layout", nil), newTicketWrapper(), mustParamsWrapper("lName", "fName", "email", "subject", "description"), methodsWrapper("POST")))
	http.Handle("/edit", adapt(nil, serveTemplateWrapper(tmpls["edit"], "layout", nil), dataWrapperOne(), mustParamsWrapper("ticket")))
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
