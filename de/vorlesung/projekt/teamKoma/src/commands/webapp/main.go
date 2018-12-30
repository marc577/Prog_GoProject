//Matrikelnummern:
//9188103
//1798794
//4717960
package main

import (
	"flag"
	"log"
	"os"
	"storagehandler"
	"webserver"
)

func main() {

	userStoreLoc := flag.String("userLoc", "../../../storage/users.json", "User Storage Path")
	ticketStoreLoc := flag.String("ticketLoc", "../../../storage/tickets/", "Ticket Storage Path")
	webPort := flag.Int("port", 8443, "https Webserver Port")
	tlsCrt := flag.String("crt", "../../../keys/server.crt", "https Webserver Certificate")
	tlsKey := flag.String("key", "../../../keys/server.key", "https Webserver Keyfile")
	htmlLoc := flag.String("htmlLoc", "../../../html", "Path to html template folder")

	flag.Parse()

	startup(*userStoreLoc, *ticketStoreLoc, *webPort, *tlsCrt, *tlsKey, *htmlLoc)
}

func startup(userStoreLoc string, ticketStoreLoc string, webPort int, tlsCrt string, tlsKey string, htmlLoc string) {

	createDirIfNotExist(ticketStoreLoc)
	createDirIfNotExist(htmlLoc)
	_, existed := createUserJSONIfNotExist(userStoreLoc)

	st := storagehandler.New(userStoreLoc, ticketStoreLoc)
	if existed == false {
		st.CreateUser("admin", "admin")
	}
	wsErr := webserver.Start(webPort, tlsCrt, tlsKey, htmlLoc, st)
	if wsErr != nil {
		log.Fatal("WebServer Error", wsErr)
	}
}

func createUserJSONIfNotExist(file string) (success bool, existed bool) {
	success = false
	existed = true
	if _, err := os.Stat(file); os.IsNotExist(err) {
		newFile, err := os.Create(file)
		existed = false
		if err != nil {
			success = false
			log.Fatal("Could not create users.json "+file+": ", err)
			return false, success
		}
		newFile.Close()
		success = true
	}

	return success, existed
}

func createDirIfNotExist(dir string) (success bool) {
	success = false
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			success = false
			log.Fatal("Could not create Folder "+dir+": ", err)
			return success
		}
		success = true
	}

	return success

}
