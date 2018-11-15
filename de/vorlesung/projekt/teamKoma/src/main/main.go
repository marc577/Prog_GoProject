package main

import (
	"flag"
	"log"
	"logging"
	"runtime"
	"strconv"
	"strings"
	"webserver"
)

func init() {
	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	logLoc := flag.String("logLoc", "../../log", "Logfile Location")
	storeLoc := flag.String("storeLoc", "../../storage", "Ticketsystem Storage Path")
	webPort := flag.Int("port", 8443, "https Webserver Port")
	tlsCrt := flag.String("crt", "../../keys/server.crt", "https Webserver Certificate")
	tlsKey := flag.String("key", "../../keys/server.key", "https Webserver Keyfile")
	htmlLoc := flag.String("htmlLoc", "../../html", "Path to html template folder")

	flag.Parse()
	logging.LogInit(*logLoc)
	logging.Info.Println(strings.Join([]string{"Flags parsed: LogLoc:", *logLoc}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: StoreLoc:", *storeLoc}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: Port:", strconv.Itoa(*webPort)}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: CRT File:", *tlsCrt}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: KEY File:", *tlsKey}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: htmlLoc:", *htmlLoc}, ""))

	wsErr := webserver.Start(*webPort, *tlsCrt, *tlsKey, *htmlLoc)
	if wsErr != nil {
		log.Fatal("WebServer Error", wsErr)
		logging.Error.Fatal("WebServer Error", wsErr)
	}
	logging.ShutdownLogging()

}
