package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/alecthomas/template"
)

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	f, fileErr := os.OpenFile("../../log/main.go.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		log.Println(fileErr)
	}
	defer f.Close()

	logger := log.New(f, "Main.go ", log.LstdFlags)
	webPort := flag.Int("port", 8443, "https Webserver Port")

	flag.Parse()
	logger.Println(joinStr("\n Flags parsed: Port:", strconv.Itoa(*webPort)))

	//http Route Handles
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/login", serveLogin)
	httpErr := http.ListenAndServeTLS(joinStr(":", strconv.Itoa(*webPort)), "../../keys/server.crt", "../../keys/server.key", nil)
	if httpErr != nil {
		logger.Println("Fatal error")
		log.Fatal("ListenAndServe: ", httpErr)
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
