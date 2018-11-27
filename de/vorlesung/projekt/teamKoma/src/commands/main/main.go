package main

import (
	"flag"
	"log"
	"logging"
	"os"
	"storagehandler"
	"strconv"
	"strings"
	"webserver"
)

func main() {

	logLoc := flag.String("logLoc", "../../../log", "Logfile Location")
	userStoreLoc := flag.String("userLoc", "../../../storage/users.json", "User Storage Path")
	ticketStoreLoc := flag.String("ticketLoc", "../../../storage/tickets/", "Ticket Storage Path")
	webPort := flag.Int("port", 8443, "https Webserver Port")
	tlsCrt := flag.String("crt", "../../../keys/server.crt", "https Webserver Certificate")
	tlsKey := flag.String("key", "../../../keys/server.key", "https Webserver Keyfile")
	htmlLoc := flag.String("htmlLoc", "../../../html", "Path to html template folder")

	flag.Parse()
	logging.LogInit(*logLoc)
	logging.Info.Println(strings.Join([]string{"Flags parsed: LogLoc:", *logLoc}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: userStoreLoc:", *userStoreLoc}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: TicketStoreLoc:", *ticketStoreLoc}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: Port:", strconv.Itoa(*webPort)}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: CRT File:", *tlsCrt}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: KEY File:", *tlsKey}, ""))
	logging.Info.Println(strings.Join([]string{"Flags parsed: htmlLoc:", *htmlLoc}, ""))

	createDirIfNotExist(logLoc)
	createDirIfNotExist(ticketStoreLoc)
	createDirIfNotExist(htmlLoc)

	st := storagehandler.New(*userStoreLoc, *ticketStoreLoc)

	wsErr := webserver.Start(*webPort, *tlsCrt, *tlsKey, *htmlLoc, st)
	if wsErr != nil {
		log.Fatal("WebServer Error", wsErr)
		logging.Error.Fatal("WebServer Error", wsErr)
	}
	logging.ShutdownLogging()

}

func createDirIfNotExist(dir string) (success bool) {
	success = false
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			success = false
			logging.Error.Panic("Could not create Folder"+dir+": ", err)
			return success
		}
		success = true
	}

	return success

}
