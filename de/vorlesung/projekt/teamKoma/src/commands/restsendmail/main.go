//Matrikelnummern:
//9188103
//1798794
//4717960
package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"storagehandler"
	"strconv"
	"strings"
)

func main() {
	//set Flags
	host := flag.String("host", "127.0.0.1", "Ticketsystem Hostname")
	port := flag.Int("port", 8443, "Ticksetsystem Webserver Port")
	user := flag.String("user", "Werner", "Your Ticketsystem Username")
	pass := flag.String("password", "password", "Your Ticketsystem Password")

	flag.Parse()

	//starts user interaction process
	interact(*host, *port, *user, *pass)
}

//interact lets the user decide which mails to be send
func interact(host string, port int, user string, pass string) {
	log.Println("Flags parsed: Host:" + host)
	log.Println("Flags Parsed: Port:" + strconv.Itoa(port))
	log.Println("Flags parsed: User:" + user)
	log.Println("Flags parsed: Password:" + pass)

	fmt.Println("Client configured to connect to " + host + ":" + strconv.Itoa(port))
	fmt.Println("Mails to send:")

	mails2send := grabMailsToSend(host, port, user, pass)
	i := 0
	for _, item := range mails2send {
		fmt.Println("ID: " + strconv.Itoa(i))
		fmt.Println("TicketID: " + item.TicketID)
		fmt.Println("CreationDate: " + item.TicketItem.CreationDate.String())
		fmt.Println("Creator: " + item.TicketItem.Creator)
		fmt.Println("Mail To: " + item.TicketItem.EmailTo)
		fmt.Println("Text to send: " + item.TicketItem.Text + "\n")
		i = i + 1
	}
	if len(mails2send) > 0 {
		cmd := "0"
		mailqueue := make([]storagehandler.Email, 0)
		idsDone := make([]int, 0)
		fmt.Println("Which IDs should be marked as sent ('all' to mark all, 'send' to finish marking, 'exit' to abort):")
		reader := bufio.NewReader(os.Stdin)
		for ok := true; ok; {
			cmd, _ = reader.ReadString('\n')
			cmd = strings.Replace(cmd, "\r\n", "", -1)
			cmd = strings.Replace(cmd, "\n", "", -1)
			if strings.Compare("all", cmd) == 0 {
				setSentFlag(host, port, user, pass, mails2send)
				ok = false
			} else {
				if strings.Compare("send", cmd) == 0 {
					setSentFlag(host, port, user, pass, mailqueue)
					ok = false
				} else {
					if strings.Compare("exit", cmd) == 0 {
						fmt.Println("bye bye")
						ok = false
					} else {
						var cmdI, err = strconv.Atoi(cmd)
						if err != nil {
							fmt.Println("Please enter correct command")
						} else {
							if cmdI >= 0 && cmdI < len(mails2send) {
								alreadyAdded := false
								for _, id := range idsDone {
									if cmdI == id {
										alreadyAdded = true
										break
									}
								}
								if !alreadyAdded {
									mailqueue = append(mailqueue, mails2send[cmdI])
									idsDone = append(idsDone, cmdI)
									fmt.Println("Added: " + mailqueue[len(mailqueue)-1].TicketID)
								} else {
									fmt.Println("ID already queued")
								}
							}
						}
						fmt.Println("What to do next?")
					}
				}
			}
		}
	} else {
		fmt.Println("No Mails to send")
	}

}

//grabMailsToSend from the API based on user auth
func grabMailsToSend(host string, port int, user string, pass string) []storagehandler.Email {
	//ignoring cert Authority Error
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//configurung http client for Request
	client := &http.Client{Transport: transCfg}
	req, err := http.NewRequest("GET", "https://"+host+":"+strconv.Itoa(port)+"/api/mail", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(user, pass)
	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Could not establish Connection. Is the server running? Hostname and port correct?", err)
	}
	tickets := make([]storagehandler.Email, 0)
	json.Unmarshal(body, &tickets)
	return tickets
}

//setSentFlag in all mails that have been marked to be send
func setSentFlag(host string, port int, user string, pass string, mails2send []storagehandler.Email) error {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transCfg}
	mails2sendJSON, err := json.Marshal(mails2send)
	if err != nil {
		log.Fatal("Could not generate correct JSON", err)
		return err
	} else {
		req, err := http.NewRequest("POST", "https://"+host+":"+strconv.Itoa(port)+"/api/mail", bytes.NewBuffer(mails2sendJSON))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			log.Fatal(err)
			return err
		}
		req.SetBasicAuth(user, pass)
		res, err := client.Do(req)
		if err != nil {
			log.Fatal("Could not establish Connection. Is the server running? Hostname and port correct?", err)
			return err
		}
		if res.StatusCode == http.StatusOK {
			log.Println("Set Sentflag in Mails")
		}
	}
	return nil

}
