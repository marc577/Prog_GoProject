package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func init(){
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	f, fileErr := os.OpenFile("log/main.go.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		log.Println(fileErr)
	}
	defer f.Close()

	logger := log.New(f, "Main.go ", log.LstdFlags)
	webPort := flag.Int("port",443,"https Webserver Port")

	flag.Parse()
	logger.Println(joinStr("\nFlags parsed: Port:", strconv.Itoa(*webPort)))

	//http Route Handles
	http.HandleFunc("/hello", HelloServer)
	httpErr := http.ListenAndServeTLS(joinStr(":", strconv.Itoa(*webPort)), "keys/server.crt", "keys/server.key", nil)
	if httpErr != nil {
		logger.Println("Fatal error")
		log.Fatal("ListenAndServe: ", httpErr)
	}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func joinStr(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}