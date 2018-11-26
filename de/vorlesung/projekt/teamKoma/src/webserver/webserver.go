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

func saveParamsWrapper(params ...string) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			ctx := r.Context()
			for _, param := range params {
				if len(r.Form.Get(param)) != 0 {
					ctx = context.WithValue(ctx, contextKey(param), r.Form.Get(param))
				}
			}
			if h != nil {
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

func redirectWrapper(path string) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := "https://" + r.Host + path
			http.Redirect(w, r, url, http.StatusFound)
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

func dataWrapper(key string, f func(s string) interface{}) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxVal := r.Context().Value(contextKey(key))
			if ctxVal != nil {
				user := ctxVal.(string)
				data := f(user)
				ctx := context.WithValue(r.Context(), contextKey("data"), data)
				if h != nil {
					h.ServeHTTP(w, r.WithContext(ctx))
				}
			} else {
				http.Error(w, http.StatusText(http.StatusNotFound)+"|"+key, http.StatusNotFound)
			}
		})
	}
}

func functionCtxWrapper(f func(w http.ResponseWriter, r *http.Request) context.Context) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := f(w, r)
			if h != nil {
				if ctx != nil {
					h.ServeHTTP(w, r.WithContext(ctx))
				} else {
					h.ServeHTTP(w, r)
				}
			}
		})
	}
}

type webContext struct {
	Data interface{}
	Path string
	User interface{}
}

func serveTemplateWrapper(t *template.Template, name string, data interface{}) adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			d := data
			if d == nil {
				ctxKey := contextKey("data")
				ctxVal := r.Context().Value(ctxKey)
				d = ctxVal
			}
			user := ""
			ctxVal := r.Context().Value(contextKey("user"))
			if ctxVal != nil {
				user = ctxVal.(string)
			}

			t.Funcs(template.FuncMap{
				"getUser": func() string { return string(user) },
			})
			path := r.URL.Path
			web := webContext{d, path, user}
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
func Start(port int, serverCertPath string, serverKeyPath string, rootPath string, st *storagehandler.StorageHandler) error {

	htmlRoot := rootPath
	defaultOpenT := template.New("").Funcs(map[string]interface{}{
		"getUser": func() string { return "" },
	})
	defaultEditT := template.New("").Funcs(map[string]interface{}{
		"getUser": func() string { return "" },
		"getAllTByProcessor": func() []string {
			return []string{"Werner"}
		},
		"getNonHolydaier": func() *[]storagehandler.User {
			user := st.GetAvailableUsers()
			return user
		},
	})

	// static files
	staticFilePath := htmlRoot + "/" + "assets"
	fs := http.FileServer(http.Dir(staticFilePath))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// templates
	tmpls := make(map[string]*template.Template)

	tmpls["index"] = template.Must(template.ParseFiles(rootPath+"/new.tmpl.html", rootPath+"/index.tmpl.html"))
	tmpls["open"] = template.Must(defaultOpenT.ParseFiles(rootPath+"/orow.tmpl.html", rootPath+"/dashboard.tmpl.html", rootPath+"/index.tmpl.html"))
	tmpls["assigned"] = template.Must(template.ParseFiles(rootPath+"/arow.tmpl.html", rootPath+"/dashboard.tmpl.html", rootPath+"/index.tmpl.html"))
	tmpls["all"] = template.Must(template.ParseFiles(rootPath+"/row.tmpl.html", rootPath+"/dashboard.tmpl.html", rootPath+"/index.tmpl.html"))
	tmpls["added"] = template.Must(template.ParseFiles(rootPath+"/added.tmpl.html", rootPath+"/index.tmpl.html"))
	tmpls["edit"] = template.Must(defaultEditT.ParseFiles(rootPath+"/ticket.tmpl.html", rootPath+"/index.tmpl.html"))

	auth := AuthenticatorFunc(st.VerifyUser)

	// frontend
	http.Handle("/open", adapt(nil, serveTemplateWrapper(tmpls["open"], "layout", nil), dataWrapper("user", func(s string) interface{} {
		return st.GetOpenTickets()
	}), basicAuthWrapper(auth)))
	http.Handle("/assigned", adapt(nil, serveTemplateWrapper(tmpls["assigned"], "layout", nil), dataWrapper("user", func(s string) interface{} {
		return st.GetNotClosedTicketsByProcessor(s)
	}), basicAuthWrapper(auth), methodsWrapper("GET")))
	http.Handle("/all", adapt(nil, serveTemplateWrapper(tmpls["all"], "layout", nil), dataWrapper("user", func(s string) interface{} {
		return st.GetTickets()
	}), basicAuthWrapper(auth), methodsWrapper("GET")))

	http.Handle("/new", adapt(nil, serveTemplateWrapper(tmpls["added"], "layout", nil), functionCtxWrapper(func(w http.ResponseWriter, r *http.Request) context.Context {
		if r.PostForm == nil {
			r.ParseForm()
		}
		if verifyEMail(r.Form.Get("email")) != true {
			http.Error(w, http.StatusText(http.StatusNotFound)+"|email", http.StatusNotFound)
			return nil
		}
		t, _ := st.CreateTicket(r.Form.Get("subject"), r.Form.Get("description"), r.Form.Get("fName"), r.Form.Get("email"), r.Form.Get("lName"))
		return context.WithValue(r.Context(), contextKey("data"), t)
	}), mustParamsWrapper("lName", "fName", "email", "subject", "description"), methodsWrapper("POST")))
	http.Handle("/edit", adapt(nil, serveTemplateWrapper(tmpls["edit"], "layout", nil), dataWrapper("ticket", func(s string) interface{} {
		//TODO: if POST update first
		t, err := st.GetTicketByID(s)
		if err != nil {
			return nil
		}
		return t
	}), saveParamsWrapper("ticket"), mustParamsWrapper("ticket"), basicAuthWrapper(auth), methodsWrapper("GET", "POST")))
	http.Handle("/edit/add", adapt(nil, functionCtxWrapper(func(w http.ResponseWriter, r *http.Request) context.Context {
		//TODO: add entry to ticket (perhaps inform kunde)
		return nil
	}), mustParamsWrapper("ticket", "description"), basicAuthWrapper(auth), methodsWrapper("POST")))
	http.Handle("/edit/free", adapt(nil, redirectWrapper("/assigned"), functionCtxWrapper(func(w http.ResponseWriter, r *http.Request) context.Context {
		//TODO: set ticket to free
		return nil
	}), mustParamsWrapper("ticket"), basicAuthWrapper(auth), methodsWrapper("GET")))

	http.Handle("/edit/combine", adapt(nil, functionCtxWrapper(func(w http.ResponseWriter, r *http.Request) context.Context {
		//TODO: combine ticket ticket and ticket toticket and redirect to the new ticket
		return nil
	}), mustParamsWrapper("ticket", "toticket"), basicAuthWrapper(auth), methodsWrapper("POST")))

	http.Handle("/assign", adapt(nil, redirectWrapper("/open"), functionCtxWrapper(func(w http.ResponseWriter, r *http.Request) context.Context {
		r.ParseForm()
		userid := r.Form.Get("user")
		ticketid := r.Form.Get("ticket")
		ticket, err := st.GetTicketByID(ticketid)
		if err == nil {
			ticket.SetTicketStateInProgress(userid)
		}
		return nil
	}), mustParamsWrapper("user", "ticket"), basicAuthWrapper(auth), methodsWrapper("GET")))

	http.Handle("/", adapt(nil, serveTemplateWrapper(tmpls["index"], "layout", nil)))

	// rest-api
	// insert ticket via mail
	http.Handle("/api/new", adapt(nil, mustParamsWrapper("lName", "fName", "email", "subject", "description"), basicAuthWrapper(auth), methodsWrapper("POST")))
	// mail sending
	http.Handle("/api/mail", adapt(nil, basicAuthWrapper(auth), methodsWrapper("GET", "POST")))

	portString := ":" + strconv.Itoa(port)
	httpErr := http.ListenAndServeTLS(portString, serverCertPath, serverKeyPath, nil)

	return httpErr
}
