package restsendmail

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"storagehandler"
	"strconv"
)

func main() {

	host := flag.String("host", "127.0.0.1", "Ticketsystem Hostname")
	port := flag.Int("port", 8443, "Ticksetsystem Webserver Port")
	user := flag.String("user", "Werner", "Your Ticketsystem Username")
	pass := flag.String("password", "password", "Your Ticketsystem Password")

	flag.Parse()
	log.Println("Flags parsed: Host:" + *host)
	log.Println("Flags Parsed: Port:" + strconv.Itoa(*port))
	log.Println("Flags parsed: User:" + *user)
	log.Println("Flags parsed: Password:" + *pass)

	fmt.Println("Client configured to connect to " + *host + ":" + strconv.Itoa(*port))
	fmt.Println("Commands to execute:")
	fmt.Println("Catch Mails to be sent: 1")
	fmt.Println("Exit : exit")

	mails2send := grabMailsToSend(*host, *port, *user, *pass)
	log.Println(mails2send)

	cmd := "1"
	fmt.Sscan(cmd)

	if cmd == "1" {
		setAllSentFlag(*host, *port, *user, *pass, mails2send)
	}

}

func grabMailsToSend(host string, port int, user string, pass string) []storagehandler.Email {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transCfg}
	req, err := http.NewRequest("GET", "https://"+host+":"+strconv.Itoa(port)+"/api/mail", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(user, pass)
	res, err := client.Do(req)
	log.Println(res)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	tickets := make([]storagehandler.Email, 0)
	json.Unmarshal(body, &tickets)
	return tickets
}

func setSentFlag() {

}

func setAllSentFlag(host string, port int, user string, pass string, mails2send []storagehandler.Email) {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transCfg}
	mails2sendJSON, err := json.Marshal(mails2send)
	if err != nil {
		log.Fatal("Could not generate correct JSON", err)
	} else {
		req, err := http.NewRequest("POST", "https://"+host+":"+strconv.Itoa(port)+"/api/mail", bytes.NewBuffer(mails2sendJSON))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(user, pass)
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode == http.StatusOK {
			log.Println("Set Sentflag in all Mails")
		}
	}

}
