package main

import (
	"flag"
	"fmt"
	"html/template"
	"logging"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

func init() {

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	logging.LogInit()
	webPort := flag.Int("port", 8443, "https Webserver Port")

	flag.Parse()
	logging.Info.Println(joinStr("Flags parsed: Port:", strconv.Itoa(*webPort)))

	//http Route Handles
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/login", serveLogin)
	httpErr := http.ListenAndServeTLS(joinStr(":", strconv.Itoa(*webPort)), "../../keys/server.crt", "../../keys/server.key", nil)
	if httpErr != nil {
		logging.Error.Fatal("ListenAndServe: ", httpErr)
	}
}

func serveIndex(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
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

func joinStr(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
