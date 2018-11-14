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

	logLoc := flag.String("logLoc", "../../log", "Logfile Location")
	storeLoc := flag.String("storeLoc", "../../storage", "Ticketsystem Storage Path")
	WebPort := flag.Int("port", 8443, "https Webserver Port")
	TLSCrt := flag.String("crt", "../../keys/server.crt", "https Webserver Certificate")
	TLSKey := flag.String("key", "../../keys/server.key", "https Webserver Keyfile")

	flag.Parse()
	logging.LogInit(*logLoc)
	logging.Info.Println(strings.Join([]string{"Flags parsed: LogLoc:", *logLoc}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: StoreLoc:", *storeLoc}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: Port:", strconv.Itoa(*WebPort)}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: CRT File:", *TLSCrt}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: KEY File:", *TLSKey}, ""))

	logging.ShutdownLogging()

	//http Route Handles
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/login", serveLogin)
	httpErr := http.ListenAndServeTLS(strings.Join([]string{":", strconv.Itoa(*WebPort)}, ""), *TLSCrt, *TLSKey, nil)
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
