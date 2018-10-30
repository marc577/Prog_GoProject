package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)
func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func main() {
	var webPort int = 8443
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(joinStr(":", strconv.Itoa(webPort)), "keys/server.crt", "keys/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func joinStr(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}